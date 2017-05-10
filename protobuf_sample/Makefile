.PHONY: deps generate build clean test

GOPATH=$(shell pwd)
export GOPATH
PATH=$(shell echo "$$GOPATH/bin:$$PATH")
export PATH

all: build

deps:
	go get -v -d .

generate:
	mkdir -p src/example
	protoc --go_out=src/example *.proto

build:
	go build -o main *.go

clean:
	rm -f main

test:
	go test .
