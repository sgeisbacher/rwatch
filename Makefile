gen-mocks:
	mockgen -typed -package mocks -source=utils/screen.go -destination mocks/screen.go Screen 
	mockgen -typed -package mocks -source=runner.go -destination mocks/executor.go -exclude_interfaces Runner Executor
