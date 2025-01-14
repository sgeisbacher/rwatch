package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	command, args := parseArgs(os.Args)
	fmt.Printf("command: %s, args: %v\n", command, args)
	screen := WebRTCScreen{}

	runnerDone := make(chan bool, 1)
	runner := LoopRunner{
		maxRunCount: 40,
		executor: func(name string, arg ...string) Executor {
			return &OsExecutor{exec.Command(name, arg...)}
		}}
	go runner.Run(&screen, runnerDone, command, args)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// TODO handle shutdown
	// go func() {
	// 	for sig := range c {
	// 	}
	// }()
	select {
	case <-c:
	case <-runnerDone:
	}
	fmt.Println("good bye!")
}

func parseArgs(args []string) (string, []string) {
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
