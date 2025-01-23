gen-mocks:
	mockgen -typed -package mocks -source=utils/screen.go -destination mocks/screen.go Screen 
	mockgen -typed -package mocks -source=runner.go -destination mocks/executor.go -exclude_interfaces Runner Executor

deploy-ui:
	test -n "$(HOST)"
	scp -r ui/index.html ui/conn.js $(HOST):/opt/rwatch/server/ui/
