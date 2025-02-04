package main

import (
	"fmt"
	"time"

	. "github.com/sgeisbacher/rwatch/utils"
)

type LoopRunner struct {
	maxRunCount int64
	executor    func(name string, arg ...string) Executor
	onDone      func()
}

func (r *LoopRunner) Run(screen Screen, done chan bool, commandName string, args []string) {
	go screen.InitScreen()
	var count int64
	for {
		count++
		cmd := r.executor(commandName, args...)
		output, err := cmd.CombinedOutput()
		runInfo := ExecutionInfo{
			CommandStr: fmt.Sprintf("%s %v", commandName, args),
			ExecTime:   time.Now(),
			ExecCount:  count,
			Output:     output,
			Success:    true,
		}
		if err != nil {
			runInfo.Success = false
			screen.SetError(fmt.Errorf("could not run command: %w\n", err))
		}
		if !cmd.WasSuccess() {
			runInfo.Success = false
			screen.SetError(fmt.Errorf("command exited with error: %w\n", err))
		}

		screen.SetOutput(runInfo)

		time.Sleep(2 * time.Second)
		if r.maxRunCount > 0 && count >= r.maxRunCount {
			fmt.Printf("maxRunCount reached (%d)\n", r.maxRunCount)
			if r.onDone != nil {
				r.onDone()
			}
			done <- true
			break
		}
	}
}
