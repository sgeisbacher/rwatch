package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

var (
	maxRunCount = flag.Int64("maxRunCount", 0, "how often the command should be run")
)

func main() {
	command, args := parseArgs(os.Args)
	fmt.Printf("command: %s, args: %v\n", command, args)
	screen := WebRTCScreen{}

	runnerDone := make(chan bool, 1)
	runner := LoopRunner{
		maxRunCount: *maxRunCount,
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
