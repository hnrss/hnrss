all: hnrss

hnrss: *.go
	go build -o $@ -ldflags "-X main.buildString=$(shell git describe --tags)"

hnrss_linux_amd64: *.go
	env GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.buildString=$(shell git describe --tags)"

run:
	./hnrss -bind 127.0.0.1:8080

clean:
	rm -f hnrss*

.PHONY: run clean
