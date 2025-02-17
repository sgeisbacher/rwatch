package utils

import (
	"fmt"
	"time"
)

type ExecutionInfo struct {
	CommandStr string    `json:"command"`
	ExecTime   time.Time `json:"exec_time"`
	ExecCount  int64     `json:"exec_count"`
	Success    bool      `json:"success"`
	Output     string    `json:"output"`
}

func (e ExecutionInfo) String() string {
	str := fmt.Sprintf("cmdStr: %q\n", e.CommandStr)
	str += fmt.Sprintf("execTime: %v\n", e.ExecTime)
	str += fmt.Sprintf("execCount: %d\n", e.ExecCount)
	str += fmt.Sprintf("success: %v\n", e.Success)
	str += fmt.Sprintf("output: %q\n", e.Output)
	return str
}
