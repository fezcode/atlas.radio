package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"atlas.radio/internal/api"
	"atlas.radio/internal/audio"
	"atlas.radio/internal/model"
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

type state int
const (
	stateBrowsing state = iota
	stateSearching
)

type Model struct {
	stations     []model.Station
	cursor       int
	input        textinput.Model
	player       *audio.Player
	state        state
	width, height int
	current      model.Station
	isPlaying    bool
	eqValues     []int
	err          error
	scrollOffset int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "TYPE TO SEARCH..."
	ti.Prompt = " > "
	ti.TextStyle = textStyle
	
	return Model{
		stations: []model.Station{},
		input:    ti,
		player:   audio.NewPlayer(),
		eqValues: make([]int, 20),
		state:    stateBrowsing,
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
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
			if m.cursor < m.scrollOffset { m.scrollOffset-- }
		case "down", "j":
			if m.cursor < len(m.stations)-1 { m.cursor++ }
			maxVisible := m.height - 10
			if m.cursor >= m.scrollOffset+maxVisible { m.scrollOffset++ }
		case "enter":
			if len(m.stations) > 0 {
				m.current = m.stations[m.cursor]
				m.isPlaying = true
				m.player.Play(m.current.URL)
			}
		}

	case stationsMsg:
		m.stations = msg
		m.cursor = 0
		m.scrollOffset = 0

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

	return m, nil
}

func (m Model) renderList() string {
	if len(m.stations) == 0 {
		return dimStyle.Render("SCANNING FREQUENCIES...")
	}

	var sb strings.Builder
	maxVisible := m.height - 10
	if maxVisible < 1 { maxVisible = 1 }
	
	end := m.scrollOffset + maxVisible
	if end > len(m.stations) { end = len(m.stations) }

	for i := m.scrollOffset; i < end; i++ {
		s := m.stations[i]
		name := s.Name
		if len(name) > 30 { name = name[:27] + "..." }
		
		line := fmt.Sprintf("%-30s", name)
		if i == m.cursor {
			sb.WriteString(selectedItemStyle.Render("> "+line) + "\n")
		} else {
			sb.WriteString(textStyle.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func (m Model) View() string {
	if m.width == 0 { return "Initializing..." }

	var mainContent string
	if m.state == stateSearching {
		mainContent = headerStyle.Render(" SEARCH FREQUENCY ") + "\n\n" + m.input.View() + "\n\n" + dimStyle.Render("ENTER to search, ESC to cancel")
	} else {
		listSide := m.renderList()
		
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
		
		playerSide += "\n\n" + dimStyle.Render("J/K: NAVIGATE\nENTER: TUNE IN\nS: SEARCH\nP: POWER OFF\nQ: QUIT")

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, 
			lipgloss.NewStyle().Width(m.width/2-4).Background(onyx).Render(listSide),
			lipgloss.NewStyle().Width(m.width/2-4).PaddingLeft(4).Background(onyx).Render(playerSide),
		)
	}

	return screenStyle.Width(m.width - 4).Height(m.height - 2).Render(mainContent)
}

func restlessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Background(onyx)
}
