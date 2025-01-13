package main

import "os/exec"

type OsExecutor struct{ cmd *exec.Cmd }

func (ose *OsExecutor) CombinedOutput() ([]byte, error) {
	return ose.cmd.CombinedOutput()
}

func (ose *OsExecutor) WasSuccess() bool {
	return ose.cmd.ProcessState.Success()
}
