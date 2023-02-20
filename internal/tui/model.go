package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type (
	stageMsg string

	jobMsg string

	model struct {
		spinner spinner.Model
		stage   string
		jobs    []string
		aborted bool
	}
)

func newModel() model {
	const windowSize = 10
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		stage:   "Initializing",
		spinner: s,
		jobs:    make([]string, windowSize),
	}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.aborted = true
		return m, tea.Quit
	case jobMsg:
		m.jobs = append(m.jobs[1:], string(msg))
		return m, nil
	case stageMsg:
		m.stage = string(msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	if m.aborted {
		s += "Aborted!"
	} else {
		s += m.spinner.View() + " " + m.stage
	}

	s += "\n\n"

	for _, res := range m.jobs {
		s += res + "\n"
	}

	if !m.aborted {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.aborted {
		s += "\n"
	}

	return appStyle.Render(s)
}
