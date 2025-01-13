package main

import (
	"fmt"
	"time"

	. "github.com/sgeisbacher/rwatch/utils"
)

type Executor interface {
	CombinedOutput() ([]byte, error)
	WasSuccess() bool
}
type Runner struct {
	maxRunCount int64
	executor    func(name string, arg ...string) Executor
}

func (r *Runner) Run(screen Screen, done chan bool, commandName string, args []string) {
	go screen.Init()
	var count int64
	for {
		count++
		cmd := r.executor(commandName, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			screen.SetError(fmt.Errorf("could not run command: %w\n", err))
			return
		}

		if !cmd.WasSuccess() {
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
		if count >= r.maxRunCount {
			done <- true
			break
		}
	}
}
