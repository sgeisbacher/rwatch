package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sgeisbacher/rwatch/utils"
)

const WIDTH = 120
const TIME_FORMAT = "02-Jan-06 15:04:05"

type TuiScreen struct {
	termPanel     viewport.Model
	eventLogPanel viewport.Model
	runInfo       utils.ExecutionInfo
}

var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("99"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusLabelNeutralStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#928f93")).
				Padding(0, 1).
				MarginRight(1)

	statusLabelOKStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#54ba6e")).
				Padding(0, 1).
				MarginRight(1)

	statusLabelFAILEDStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF0000")).
				Padding(0, 1)

	statusTextStyle = lipgloss.NewStyle().Inherit(statusBarStyle).Foreground(lipgloss.Color("#FFFFFF"))
)

type tickMsg time.Time

func (ts *TuiScreen) InitScreen() {
	ts.termPanel = viewport.New(WIDTH, 20)
	ts.termPanel.Style = focusedBorderStyle
}

func (ts *TuiScreen) SetOutput(info utils.ExecutionInfo) {
	ts.runInfo = info
}

func (ts *TuiScreen) SetError(err error) {

}

func (ts *TuiScreen) Init() tea.Cmd {
	return tick()
}

func (ts *TuiScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO remove duplicate
	quitKey := key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC", "to quit"),
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitKey):
			return ts, tea.Quit
		}
		// case tea.WindowSizeMsg:
		// 	m.height = msg.Height
		// 	m.width = msg.Width
	}

	// m.updateKeybindings()
	// m.sizeInputs()

	// return m, tea.Batch(cmds...)
	return ts, tick()
}

func (ts *TuiScreen) View() string {
	// TODO remove duplicate
	quitKey := key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ESC", "to quit"),
	)
	help := help.New()
	helpLine := help.ShortHelpView([]key.Binding{
		// 	m.keymap.next,
		// 	m.keymap.prev,
		// 	m.keymap.add,
		// 	m.keymap.remove,
		quitKey,
	})

	w := lipgloss.Width
	ts.termPanel.SetContent(string(ts.runInfo.Output))
	var sline_lstatus string
	if ts.runInfo.Success {
		sline_lstatus = statusLabelOKStyle.Render("OK")
	} else {
		sline_lstatus = statusLabelFAILEDStyle.Render("FAILED")
	}
	sline_rstatus := statusLabelNeutralStyle.
		MarginRight(0).
		Render("UNKNOWN")
	sline_rlabel := statusTextStyle.
		Render("WebRTC:")
	sline_runs := statusTextStyle.
		AlignHorizontal(lipgloss.Center).
		Width(WIDTH - w(sline_lstatus) - w(sline_rlabel) - w(sline_rstatus)).
		Render(fmt.Sprintf("Run: %d (%v) ", ts.runInfo.ExecCount, ts.runInfo.ExecTime.Format(TIME_FORMAT)))
	statusLine := lipgloss.JoinHorizontal(lipgloss.Top, sline_lstatus, sline_runs, sline_rlabel, sline_rstatus)

	views := []string{
		statusLine,
		ts.termPanel.View(),
		helpLine,
	}
	return lipgloss.JoinVertical(lipgloss.Top, views...)
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
