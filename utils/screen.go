package utils

type Screen interface {
	InitScreen()
	Run(runnerDone chan bool)
	SetOutput(info ExecutionInfo)
	SetError(err error)
	Done()
}
