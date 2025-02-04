package main

import (
	"fmt"

	"github.com/sgeisbacher/rwatch/utils"
)

type SimpleScreen struct{}

func (s *SimpleScreen) InitScreen() {

}

func (s *SimpleScreen) SetOutput(info utils.ExecutionInfo) {
	fmt.Printf("--\nOUTPUT (run: %d):\n\n%s\n", info.ExecCount, info.Output)
}

func (s *SimpleScreen) SetError(err error) {
}
