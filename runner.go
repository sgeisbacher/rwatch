package main

import (
	"fmt"
	"os/exec"
	"time"
)

type ExecutionInfo struct {
	CommandStr string
	ExecTime   time.Time
	ExecCount  int64
	Success    bool
	Output     []byte
}

func run(screen Screen, done chan bool, commandName string, args []string) {
	go screen.Init()
	var count int64
	for {
		count++
		fmt.Printf("run: %d\n", count)
		cmd := exec.Command(commandName, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			screen.SetError(fmt.Errorf("could not run command: %w\n", err))
			return
		}

		if !cmd.ProcessState.Success() {
			screen.SetError(fmt.Errorf("command exited with error: %w\n", err))
			return
		}

		screen.SetOutput(ExecutionInfo{
			CommandStr: fmt.Sprintf("%s %v", commandName, args),
			ExecTime:   time.Now(),
			ExecCount:  count,
			Output:     output,
		})
		time.Sleep(2 * time.Second)
		if count >= 40 {
			done <- true
			break
		}
	}
}
