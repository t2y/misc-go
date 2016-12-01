.PHONY: build clean test

GOPATH=$(shell pwd)
export GOPATH

build:
	go build -o main main.go

clean:
	rm -f main

test:
	go test .
