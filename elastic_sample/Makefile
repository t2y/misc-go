.PHONY: build clean test

GOPATH=$(shell pwd)
export GOPATH

build:
	go build -o bin/main myapp

clean:
	rm -f bin/main

test:
	go test myapp/...
