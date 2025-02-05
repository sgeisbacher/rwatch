test:
	go test -v ./...

system-test:
	go test -tags system_test -v ./...

gen-mocks:
	mockgen -typed -package mocks -source=utils/screen.go -destination mocks/screen.go Screen 
	mockgen -typed -package mocks -source=runner.go -destination mocks/executor.go -exclude_interfaces Runner Executor

deploy-ui:
	test -n "$(HOST)"
	cd ui && npm run build
	ssh $(HOST) "rm -rf $(HOST):/opt/rwatch/server/ui/"
	scp -r ui/build/* $(HOST):/opt/rwatch/server/ui/
