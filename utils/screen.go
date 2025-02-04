package utils

type Screen interface {
	InitScreen()
	SetOutput(info ExecutionInfo)
	SetError(err error)
}
