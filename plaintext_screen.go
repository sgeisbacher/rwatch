package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sgeisbacher/rwatch/utils"
)

type PlainTextScreen struct {
	appState *appStateManager
}

func (s *PlainTextScreen) InitScreen() {
	go printSessionId(s.appState)
}

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

// TODO rewrite to eventlistener
func printSessionId(appState *appStateManager) {
	for {
		if appState.GetWebRTCSessionId() != "" {
			fmt.Printf("Session-ID: %s\n", appState.GetWebRTCSessionId())
			fmt.Printf("Session-URL: %s\n", appState.GenSessionUrl("/"))
			break
		}
		time.Sleep(time.Second)
	}
}

func (s *PlainTextScreen) SetError(err error) {}

func (s *PlainTextScreen) Done() {}
