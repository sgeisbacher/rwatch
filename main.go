package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/sgeisbacher/rwatch/utils"
)

var (
	maxRunCount = flag.Int64("maxRunCount", 0, "how often the command should be run")
)

func main() {
	command, args := parseArgs(os.Args)
	fmt.Printf("command: %s, args: %v\n", command, args)

	// Setup WebRTC Screen
	webRTCScreen := WebRTCScreen{}

	// Setup WebRTC Screen
	tuiScreen := SimpleScreen{} // TuiScreen{}

	// setup bubbletea
	// tui := tea.NewProgram(&tuiScreen, tea.WithAltScreen())

	// Setup Runner
	runnerDone := make(chan bool, 1)
	runner := LoopRunner{
		maxRunCount: *maxRunCount,
		executor: func(name string, arg ...string) Executor {
			return &OsExecutor{exec.Command(name, arg...)}
		},
		//onDone: func() { tui.Quit() },
	}
	screen := MultiplexerScreen{
		screens: []utils.Screen{&tuiScreen, &webRTCScreen},
	}
	go runner.Run(&screen, runnerDone, command, args)

	// start tui
	// if _, err := tui.Run(); err != nil {
	// 	fmt.Printf("E: bubbletea: %v\n", err)
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// TODO handle shutdown
	// go func() {
	// 	for sig := range c {
	// 	}
	// }()
	select {
	case <-c:
		fmt.Println("got TERM signal")
	case <-runnerDone:
		fmt.Println("runner done")
	}
	fmt.Println("good bye!")
}

func parseArgs(args []string) (string, []string) {
	// args
	flag.Parse()

	// command
	position := -1
	for i, arg := range args {
		if arg == "--" {
			position = i + 1
			break
		}
	}
	if position == -1 {
		panic("you need to have -- as command separator")
	}
	return args[position], args[position+1:]
}
