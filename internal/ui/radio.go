package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"atlas.radio/internal/api"
	"atlas.radio/internal/audio"
	"atlas.radio/internal/model"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	amber  = lipgloss.Color("#FFB642")
	onyx   = lipgloss.Color("#050505")
	rusty  = lipgloss.Color("#5E4737")

	screenStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(amber).
			Background(onyx)

	headerStyle = lipgloss.NewStyle().
			Foreground(onyx).
			Background(amber).
			Padding(0, 1).
			Bold(true)

	textStyle = lipgloss.NewStyle().Foreground(amber).Background(onyx)
	dimStyle  = lipgloss.NewStyle().Foreground(rusty).Background(onyx)
	selectedItemStyle = lipgloss.NewStyle().Foreground(onyx).Background(amber)
)

type item struct {
	station model.Station
}
func (i item) Title() string       { return i.station.Name }
func (i item) Description() string { return fmt.Sprintf("%s | %s", i.station.Country, i.station.Tags) }
func (i item) FilterValue() string { return i.station.Name }

type state int
const (
	stateBrowsing state = iota
	stateSearching
)

type Model struct {
	list         list.Model
	input        textinput.Model
	player       *audio.Player
	state        state
	width, height int
	current      model.Station
	isPlaying    bool
	eqValues     []int
	err          error
}

func NewModel() Model {
	// Custom delegate to ensure Onyx background everywhere
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = textStyle
	d.Styles.SelectedTitle = selectedItemStyle
	d.Styles.NormalDesc = dimStyle
	d.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(onyx).Background(rusty)

	l := list.New([]list.Item{}, d, 0, 0)
	l.Title = " PIP-BOY 3000 RADIO "
	l.Styles.Title = headerStyle
	l.Styles.HelpStyle = dimStyle.PaddingLeft(2)
	l.SetShowStatusBar(false)
	
	ti := textinput.New()
	ti.Placeholder = "Enter station name or tag..."
	ti.TextStyle = textStyle
	
	return Model{
		list:   l,
		input:  ti,
		player: audio.NewPlayer(),
		eqValues: make([]int, 20),
		state: stateBrowsing,
	}
}

func (m Model) Stop() {
	if m.player != nil {
		m.player.Stop()
	}
}

type stationsMsg []model.Station
type tickMsg time.Time

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.searchCmd(""), tick())
}

func (m Model) searchCmd(q string) tea.Cmd {
	return func() tea.Msg {
		stations, err := api.SearchStations(q)
		if err != nil {
			return err
		}
		return stationsMsg(stations)
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Account for padding and border
		m.list.SetSize(msg.Width/2-10, msg.Height-10)

	case tea.KeyMsg:
		if m.state == stateSearching {
			switch msg.String() {
			case "enter":
				m.state = stateBrowsing
				return m, m.searchCmd(m.input.Value())
			case "esc":
				m.state = stateBrowsing
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		case "s", "/":
			m.state = stateSearching
			m.input.Focus()
			return m, nil
		case "p":
			m.isPlaying = false
			m.player.Stop()
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.current = i.station
				m.isPlaying = true
				m.player.Play(m.current.URL)
			}
		}

	case stationsMsg:
		items := make([]list.Item, len(msg))
		for i, s := range msg {
			items[i] = item{station: s}
		}
		m.list.SetItems(items)

	case error:
		m.err = msg

	case tickMsg:
		if m.isPlaying {
			for i := range m.eqValues {
				m.eqValues[i] = rand.Intn(5)
			}
		}
		return m, tick()
	}

	if m.state == stateBrowsing {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 { return "Initializing..." }

	var mainContent string
	if m.state == stateSearching {
		mainContent = headerStyle.Render(" SEARCH FREQUENCY ") + "\n\n" + m.input.View() + "\n\n" + dimStyle.Render("ENTER to search, ESC to cancel")
	} else {
		listSide := m.list.View()
		
		playerSide := "\n" + headerStyle.Render(" SIGNAL DATA ") + "\n"
		if m.err != nil {
			playerSide += restlessStyle().Render("ERROR: OFFLINE") + "\n"
		} else if m.isPlaying {
			playerSide += textStyle.Render("STATION: " + m.current.Name) + "\n"
			playerSide += dimStyle.Render("SIGNAL: " + m.current.Country) + "\n\n"
			
			playerSide += textStyle.Render("FREQ VISUALIZER:") + "\n"
			var hEq strings.Builder
			for h := 4; h >= 0; h-- {
				for _, v := range m.eqValues {
					if v > h { hEq.WriteString("█") } else { hEq.WriteString(" ") }
				}
				hEq.WriteString("\n")
			}
			playerSide += textStyle.Render(hEq.String())
		} else {
			playerSide += dimStyle.Render("RECEIVER IDLE\nNO SIGNAL DETECTED") + "\n"
		}
		
		playerSide += "\n\n" + dimStyle.Render("ENTER: TUNE IN\nS: SEARCH\nP: POWER OFF\nQ: QUIT")

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, 
			lipgloss.NewStyle().Width(m.width/2-4).Render(listSide),
			lipgloss.NewStyle().Width(m.width/2-4).PaddingLeft(4).Render(playerSide),
		)
	}

	// Apply global Pip-Boy screen style
	return screenStyle.Width(m.width - 4).Height(m.height - 2).Render(mainContent)
}

func restlessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Background(onyx)
}
