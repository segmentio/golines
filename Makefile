export GO111MODULE=on

.PHONY: all
all: build test

.PHONY: build
build:
	goreleaser build --clean --snapshot --single-target

.PHONY: test
test: vet
	go test -count=1 -cover -coverprofile=coverage.out .
	go tool cover -func=coverage.out

.PHONY: vet
vet:
	go vet .

.PHONY: regenerate
regenerate:
	REGENERATE_TEST_OUTPUTS=true go test .

.PHONY: graph
graph:
	cd generate && go run .

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: format
format:
	goimports -w ./*go
