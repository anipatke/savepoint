.PHONY: build test run clean build-linux build-darwin build-windows build-all build-npm dist verify-dist package-check smoke-test ci install-hooks

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

build-windows:
	go run ./internal/buildtool -version "$(VERSION)" build-windows

build-all: build-linux build-darwin build-windows

build-npm:
	go run ./internal/buildtool -version "$(VERSION)" build-npm

dist:
	go run ./internal/buildtool -version "$(VERSION)" dist

verify-dist:
	go run ./internal/buildtool -version "$(VERSION)" verify-dist

package-check: build-npm
	npm test
	npm pack --dry-run --ignore-scripts

smoke-test:
	go run ./internal/buildtool -version "$(VERSION)" smoke-test

ci: test build dist package-check

install-hooks:
	git config core.hooksPath scripts/git-hooks
