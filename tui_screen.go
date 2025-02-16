package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sgeisbacher/rwatch/utils"
)

const WIDTH = 120
const TIME_FORMAT = "15:04:05" //"02-Jan-06 15:04:05"

type TuiScreen struct {
	appState *appStateManager
	model    *TuiBubbleTeaModel
	tui      *tea.Program
}
type TuiBubbleTeaModel struct {
	appState      *appStateManager
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
	ts.model = &TuiBubbleTeaModel{
		appState: ts.appState,
	}
	ts.tui = tea.NewProgram(ts.model, tea.WithAltScreen())
}

func (ts *TuiScreen) Run(runnerDone chan bool) {
	for {
		if ts.tui != nil {
			break
		}
	}
	if _, err := ts.tui.Run(); err != nil {
		fmt.Printf("E: bubbletea: %v\n", err)
	}
}

func (ts *TuiScreen) SetOutput(info utils.ExecutionInfo) {
	if ts.model != nil {
		ts.model.runInfo = info
	}
}

func (ts *TuiScreen) SetError(err error) {}

func (ts *TuiScreen) Done() {
	ts.tui.Quit()
}

func (ts *TuiBubbleTeaModel) Init() tea.Cmd {
	ts.termPanel = viewport.New(WIDTH, 20)
	ts.termPanel.Style = focusedBorderStyle
	ts.eventLogPanel = viewport.New(WIDTH, 10)
	ts.eventLogPanel.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#817e82")).
		MarginLeft(2)
	return tick()
}

func (ts *TuiBubbleTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (tm *TuiBubbleTeaModel) View() string {
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

	// logs
	logMessages := []string{}
	for _, logEvent := range tm.appState.logs {
		msg := fmt.Sprintf("%s %s", logEvent.timestamp.Format(TIME_FORMAT), logEvent.msg)
		logMessages = append(logMessages, msg)
	}
	tm.eventLogPanel.SetContent(strings.Join(logMessages, "\n"))
	tm.eventLogPanel.SetYOffset(len(logMessages))

	// term
	tm.termPanel.SetContent(string(tm.runInfo.Output))

	// status line
	var sline_lstatus string
	if tm.runInfo.Success {
		sline_lstatus = statusLabelOKStyle.Render("OK")
	} else {
		sline_lstatus = statusLabelFAILEDStyle.Render("FAILED")
	}
	sline_rstatus := statusLabelNeutralStyle.
		MarginRight(0).
		Render(tm.appState.Current().Name)
	sline_rlabel := statusTextStyle.
		Render("WebRTC:")
	sline_runs := statusTextStyle.
		AlignHorizontal(lipgloss.Center).
		Width(WIDTH - w(sline_lstatus) - w(sline_rlabel) - w(sline_rstatus)).
		Render(fmt.Sprintf("Run: %d (%v) ", tm.runInfo.ExecCount, tm.runInfo.ExecTime.Format(TIME_FORMAT)))
	statusLine := lipgloss.JoinHorizontal(lipgloss.Top, sline_lstatus, sline_runs, sline_rlabel, sline_rstatus)

	var link string
	if tm.appState.GetWebRTCSessionId() != "" {
		link = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Render(tm.appState.GenSessionUrl("/"))
	}
	footerFiller := lipgloss.NewStyle().Width(WIDTH - w(helpLine) - w(link)).Render("")
	footerLine := lipgloss.JoinHorizontal(lipgloss.Top, helpLine, footerFiller, link)

	// join
	views := []string{
		statusLine,
		tm.termPanel.View(),
		tm.eventLogPanel.View(),
		"",
		footerLine,
	}
	return lipgloss.JoinVertical(lipgloss.Top, views...)
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
