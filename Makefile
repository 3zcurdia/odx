.PHONY: build install clean

BINARY_NAME=odx
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} cmd/odx/main.go

install:
	go install ${LDFLAGS} ./cmd/odx

clean:
	rm -rf bin/
	go clean

test:
	go test ./...
