.PHONY: build clean test

DEFAULT_GO_PATH := $(HOME)/go
AWS_SDK_GO_VERSION := 1.29.17
AWS_SDK_GO_REPOSITORY := $(DEFAULT_GO_PATH)/pkg/mod/github.com/aws/aws-sdk-go@v$(AWS_SDK_GO_VERSION)
PACKGE_NAME := main

mockS3.go:
	mockgen -source $(AWS_SDK_GO_REPOSITORY)/service/s3/s3iface/interface.go -package $(PACKGE_NAME) -destination ./$@

build:
	go build -o $(PACKGE_NAME) $(PACKGE_NAME).go

clean:
	rm -f $(PACKGE_NAME) mockS3.go

test: mockS3.go
	go test -v .
