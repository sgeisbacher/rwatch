package main

type Screen interface {
	Init()
	SetOutput(info ExecutionInfo)
	SetError(err error)
}
