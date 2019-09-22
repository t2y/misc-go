.PHONY: mod build clean test

all: build

mod:
	go mod init github.com/t2y/misc/validator_sample

build:
	go build -o main main.go

clean:
	rm -f main

test:
	go test .
