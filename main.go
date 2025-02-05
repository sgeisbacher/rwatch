package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/sgeisbacher/rwatch/utils"
)

var (
	maxRunCount        = flag.Int64("max-run-count", 0, "how often the command should be run")
	usePlainTextScreen = flag.Bool("plain-text-screen", false, "dont show command output in fancy bubbletea-screen but just simple plaintext-screen")
)

func main() {
	command, args := parseArgs(os.Args)
	fmt.Printf("command: %s, args: %v\n", command, args)

	// Setup AppState
	appState := createAppState()

	// Setup local screen
	var localScreen utils.Screen = &PlainTextScreen{}
	if !*usePlainTextScreen {
		localScreen = &TuiScreen{appState: appState}
	}

	// Setup WebRTC Screen
	webRTCScreen := &WebRTCScreen{appState: appState}

	// Setup Runner
	runnerDone := make(chan bool, 1)
	runner := LoopRunner{
		maxRunCount: *maxRunCount,
		executor: func(name string, arg ...string) Executor {
			return &OsExecutor{exec.Command(name, arg...)}
		},
	}
	screen := MultiplexerScreen{
		screens: []utils.Screen{localScreen, webRTCScreen},
	}
	go runner.Run(&screen, runnerDone, command, args)

	localScreen.Run(runnerDone)
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
