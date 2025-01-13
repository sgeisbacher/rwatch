package utils

type Screen interface {
	Init()
	SetOutput(info ExecutionInfo)
	SetError(err error)
}
