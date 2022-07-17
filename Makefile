.PHONY: pre-commit test generate build run
.DEFAULT_GOAL := test

PROG := $(shell basename $(PWD))

pre-commit:
	go mod tidy
	go vet

test:
	@go run github.com/onsi/ginkgo/v2/ginkgo -r --procs=1 --compilers=1 --randomize-all --randomize-suites --fail-on-pending --keep-going --cover --coverprofile=cover.profile.out --race --trace --json-report=report.json --timeout=15s

generate:
	@go generate -mod=mod ./...

build:
	@go build -mod=mod

run:
	./$(PROG)
