package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sgeisbacher/rwatch/utils"
)

type PlainTextScreen struct{}

func (s *PlainTextScreen) InitScreen() {}

func (s *PlainTextScreen) Run(runnerDone chan bool) {
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
}

func (s *PlainTextScreen) SetOutput(info utils.ExecutionInfo) {
	fmt.Printf("--\nOUTPUT (run: %d):\n\n%s\n", info.ExecCount, info.Output)
}

func (s *PlainTextScreen) SetError(err error) {}

func (s *PlainTextScreen) Done() {}
