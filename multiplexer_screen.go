package main

import "github.com/sgeisbacher/rwatch/utils"

type MultiplexerScreen struct {
	screens []utils.Screen
}

func (ms *MultiplexerScreen) InitScreen() {
	for _, screen := range ms.screens {
		screen.InitScreen()
	}
}

func (ms *MultiplexerScreen) SetOutput(info utils.ExecutionInfo) {
	for _, screen := range ms.screens {
		screen.SetOutput(info)
	}
}

func (ms *MultiplexerScreen) SetError(err error) {
	for _, screen := range ms.screens {
		screen.SetError(err)
	}
}
