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
	// Pip-Boy Amber
	amber = lipgloss.Color("#FFB642")
	dim   = lipgloss.Color("#5E4737")

	// Styles
	mainStyle = lipgloss.NewStyle().
			Foreground(amber)

	titleStyle = lipgloss.NewStyle().
			Foreground(amber).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(amber).
			Padding(0, 1).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(amber).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(amber).
			Border(lipgloss.NormalBorder()).
			BorderForeground(amber).
			Padding(0, 1).
			MarginRight(1)

	activeButtonStyle = lipgloss.NewStyle().
			Foreground(amber).
			Border(lipgloss.ThickBorder()).
			BorderForeground(amber).
			Padding(0, 1).
			MarginRight(1).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(amber).
				Bold(true)

	dimTextStyle = lipgloss.NewStyle().
			Foreground(dim)
)

type state int

const (
	stateBrowsing state = iota
	stateSearching
)

type Model struct {
	stations      []model.Station
	cursor        int
	input         textinput.Model
	player        *audio.Player
	state         state
	width, height int
	current       model.Station
	isPlaying     bool
	isLoading     bool
	eqValues      []int
	err           error
	scrollOffset  int
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "FREQUENCY SCAN..."
	ti.Prompt = " > "
	ti.TextStyle = mainStyle

	return Model{
		stations:  []model.Station{},
		input:     ti,
		player:    audio.NewPlayer(),
		eqValues:  make([]int, 30),
		state:     stateBrowsing,
		isLoading: true,
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
				m.isLoading = true
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
		case "p", " ":
			if m.isPlaying {
				m.isPlaying = false
				m.player.Stop()
			} else if len(m.stations) > 0 {
				m.current = m.stations[m.cursor]
				m.isPlaying = true
				m.err = m.player.Play(m.current.URL)
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			if m.cursor < m.scrollOffset {
				m.scrollOffset = m.cursor
			}
		case "down", "j":
			if m.cursor < len(m.stations)-1 {
				m.cursor++
			}
			maxVisible := m.height - 18
			if maxVisible < 1 {
				maxVisible = 1
			}
			if m.cursor >= m.scrollOffset+maxVisible {
				m.scrollOffset++
			}
		case "enter":
			if len(m.stations) > 0 {
				m.current = m.stations[m.cursor]
				m.isPlaying = true
				m.err = m.player.Play(m.current.URL)
			}
		}

	case stationsMsg:
		m.stations = msg
		m.cursor = 0
		m.scrollOffset = 0
		m.isLoading = false
		m.err = nil

	case error:
		m.err = msg
		m.isLoading = false

	case tickMsg:
		if m.isPlaying {
			for i := range m.eqValues {
				m.eqValues[i] = rand.Intn(8)
			}
		}
		return m, tick()
	}

	return m, nil
}

func (m Model) renderHeader() string {
	return titleStyle.
		Width(m.width - 10).
		Align(lipgloss.Center).
		Render("ATLAS RADIO RECEIVER - MODEL 3000")
}

func (m Model) renderDisplay() string {
	width := m.width - 6
	if width < 40 {
		width = 40
	}

	var content string
	if m.isLoading {
		content = "SCANNING GLOBAL FREQUENCIES...\nPLEASE STAND BY."
	} else if m.err != nil {
		content = fmt.Sprintf("SIGNAL INTERFERENCE DETECTED\nERROR: %s", strings.ToUpper(m.err.Error()))
	} else if m.isPlaying {
		name := strings.ToUpper(m.current.Name)
		if len(name) > width-15 {
			name = name[:width-18] + "..."
		}
		content = fmt.Sprintf("NOW TUNED: %s\nORIGIN   : %s", name, strings.ToUpper(m.current.Country))
	} else {
		content = "RECEIVER IDLE\nNO SIGNAL DETECTED"
	}

	return borderStyle.Width(width).Render(content)
}

func (m Model) renderControls() string {
	play := buttonStyle.Render("PLAY [P]")
	if m.isPlaying {
		play = activeButtonStyle.Render("PLAYING")
	}

	stop := buttonStyle.Render("STOP [P]")
	if !m.isPlaying {
		stop = activeButtonStyle.Render("IDLE")
	}

	scan := buttonStyle.Render("SCAN [S]")
	exit := buttonStyle.Render("OFF [Q]")

	return lipgloss.JoinHorizontal(lipgloss.Center, play, stop, scan, exit)
}

func (m Model) renderList() string {
	maxVisible := m.height - 18
	if maxVisible < 1 {
		maxVisible = 1
	}

	end := m.scrollOffset + maxVisible
	if end > len(m.stations) {
		end = len(m.stations)
	}

	var sb strings.Builder
	for i := m.scrollOffset; i < end; i++ {
		s := m.stations[i]
		name := s.Name
		if len(name) > 50 {
			name = name[:47] + "..."
		}

		if i == m.cursor {
			sb.WriteString(selectedItemStyle.Render(fmt.Sprintf(" > %-50s ", name)) + "\n")
		} else {
			sb.WriteString(mainStyle.Render(fmt.Sprintf("   %-50s ", name)) + "\n")
		}
	}
	return sb.String()
}

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var mainContent string
	if m.state == stateSearching {
		mainContent = lipgloss.JoinVertical(lipgloss.Center,
			activeButtonStyle.Render(" SEARCH FREQUENCY "),
			"\n",
			m.input.View(),
			"\n",
			dimTextStyle.Render("ENTER TO SCAN / ESC TO ABORT"),
		)
	} else {
		mainContent = lipgloss.JoinVertical(lipgloss.Center,
			m.renderHeader(),
			"\n",
			m.renderDisplay(),
			"\n",
			m.renderControls(),
			"\n",
			m.renderList(),
		)
	}

	return mainStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(mainContent)
}
