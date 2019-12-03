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
test: vet
	go test .

.PHONY: vet
vet:
	go vet .

.PHONY: regenerate
regenerate:
	REGENERATE_TEST_OUTPUTS=true go test .

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: format
format:
	goimports -w ./*go
