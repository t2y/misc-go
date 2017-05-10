.PHONY: deps generate build clean test

GOPATH=$(shell pwd)
export GOPATH

path=${GOPATH}/bin:$(PATH)
export PATH := ${path}

all: build

deps:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -v -d .

generate:
	@echo $(PATH)
	mkdir -p src/example
	protoc --go_out=src/example *.proto

build:
	go build -o main *.go

clean:
	rm -f main

test:
	go test .
