hnrss_linux_amd64:
	env GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-s -w -X main.buildString=$(shell git describe --tags)"
