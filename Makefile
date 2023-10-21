export GO111MODULE=on

.PHONY: all
all: install test

.PHONY: install
install:
	go install .

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test -count=1 -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: regenerate
regenerate:
	REGENERATE_TEST_OUTPUTS=true go test .

.PHONY: graph
graph:
	go generate ./...

.PHONY: format
format:
	pre-commit run --all-files
