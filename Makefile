.PHONY: build test test-focused test-fast test-full run clean build-linux build-darwin build-windows build-all build-npm dist verify-dist package-check smoke-test ci install-hooks

VERSION ?=

build:
	go run ./internal/buildtool -version "$(VERSION)" build

# Full host-platform Go suite. JSON output feeds the package/test timing summary.
test:
	go run ./internal/buildtool test -json -count=1 ./...

# Iteration aid: make test-focused TEST='TestName' [PKGS=./package].
test-focused:
	go run ./internal/buildtool focused-test "$(TEST)" $(if $(PKGS),$(PKGS),./...)

# Ordinary Task handoff gate. T013's three expensive migration scenarios stay in full.
test-fast:
	go run ./internal/buildtool test -json -count=1 -skip '^(TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability|TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits|TestEndToEnd_goldenIsReproducible)$$' ./...

# Migration/platform-sensitive Task and Objective/Release integration gate.
test-full: test build-all

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

ci: test-full build dist package-check

install-hooks:
	git config core.hooksPath scripts/git-hooks
