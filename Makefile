.PHONY: build test run clean build-linux build-darwin build-all dist smoke-test

VERSION ?=

build:
	go run ./internal/buildtool -version "$(VERSION)" build

test:
	go test ./...

run:
	go run main.go

clean:
	go run ./internal/buildtool -version "$(VERSION)" clean

build-linux:
	go run ./internal/buildtool -version "$(VERSION)" build-linux

build-darwin:
	go run ./internal/buildtool -version "$(VERSION)" build-darwin

build-all: build-linux build-darwin

dist:
	go run ./internal/buildtool -version "$(VERSION)" dist

smoke-test:
	go run ./internal/buildtool -version "$(VERSION)" smoke-test
