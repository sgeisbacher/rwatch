package main

import (
	"fmt"
	"testing"

	mocks "github.com/sgeisbacher/rwatch/mocks"
	"github.com/sgeisbacher/rwatch/utils"

	// "github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type ExecutionInfoMatcher struct {
	cmdStrMatcher    gomock.Matcher
	execCountMatcher gomock.Matcher
	successMatcher   gomock.Matcher
	outputMatcher    gomock.Matcher
}

func (m ExecutionInfoMatcher) Matches(x any) bool {
	info := x.(utils.ExecutionInfo)
	if m.cmdStrMatcher != nil && !m.cmdStrMatcher.Matches(info.CommandStr) {
		return false
	}
	if m.execCountMatcher != nil && !m.execCountMatcher.Matches(info.ExecCount) {
		return false
	}
	if m.outputMatcher != nil && !m.outputMatcher.Matches(info.Output) {
		return false
	}
	return true
}

func (m ExecutionInfoMatcher) String() string {
	str := ""
	if m.cmdStrMatcher != nil {
		str += fmt.Sprintf("cmdStr: %s\n", m.cmdStrMatcher.String())
	}
	if m.execCountMatcher != nil {
		str += fmt.Sprintf("count: %s\n", m.execCountMatcher.String())
	}
	if m.outputMatcher != nil {
		str += fmt.Sprintf("ouput: %s\n", m.outputMatcher.String())
	}
	return str
}

func TestRunnerHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	screenMock := mocks.NewMockScreen(ctrl)
	executorMock := mocks.NewMockExecutor(ctrl)
	runner := Runner{
		maxRunCount: 1,
		executor: func(name string, arg ...string) Executor {
			return executorMock
		},
	}
	setOuputMatcher := ExecutionInfoMatcher{
		cmdStrMatcher:    gomock.Eq("ls [-l /tmp]"),
		execCountMatcher: gomock.Eq(int64(1)),
		outputMatcher:    gomock.Eq([]byte("-rw-r--r-- 1 stefan staff 5271 Jan  13 11:18 data.txt")),
	}

	screenMock.EXPECT().Init().Times(1)
	screenMock.EXPECT().SetOutput(setOuputMatcher).Times(1)
	executorMock.EXPECT().CombinedOutput().Times(1).Return([]byte("-rw-r--r-- 1 stefan staff 5271 Jan  13 11:18 data.txt"), nil)
	executorMock.EXPECT().WasSuccess().Times(1).Return(true)

	done := make(chan bool, 1)
	runner.Run(screenMock, done, "ls", []string{"-l", "/tmp"})
}
