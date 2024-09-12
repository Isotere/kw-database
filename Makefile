SHELL   := /bin/bash -euo pipefail
TIMEOUT := 10s
GOFLAGS := -mod=vendor
BINPATH = $(PWD)/bin
TIMEOUT_UNIT_TESTS := 10s
TIMEOUT_INT_TESTS := 5m

.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: test
test:
	@go test -race -timeout $(TIMEOUT_UNIT_TESTS) ./...

.PHONY: test-integration
test-integration:
	@go test -race -timeout $(TIMEOUT_INT_TESTS) -tags integration -p 1 ./tests/...

.PHONY: test-all
test-all:
	@go test -race -timeout $(TIMEOUT_INT_TESTS) -tags integration -p 1 ./...

.PHONY: generate
generate:
	@go generate ./...

.PHONY: lint
lint:
	@golangci-lint run --out-format colored-line-number
