package utils

import (
	"fmt"
	"time"
)

type ExecutionInfo struct {
	CommandStr string
	ExecTime   time.Time
	ExecCount  int64
	Success    bool
	Output     []byte
}

func (e ExecutionInfo) String() string {
	str := fmt.Sprintf("cmdStr: %s\n", e.CommandStr)
	str += fmt.Sprintf("execTime: %v\n", e.ExecTime)
	str += fmt.Sprintf("execCount: %d\n", e.ExecCount)
	str += fmt.Sprintf("success: %v\n", e.Success)
	str += fmt.Sprintf("output: %s\n", e.Output)
	return str
}
