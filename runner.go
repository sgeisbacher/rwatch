package main

import (
	"github.com/sgeisbacher/rwatch/utils"
)

type Executor interface {
	CombinedOutput() ([]byte, error)
	WasSuccess() bool
}
type Runner interface {
	Run(screen utils.Screen, done chan bool, commandName string, args []string)
}
