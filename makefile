PACKAGES=$(shell go list ./... | grep -v 'tests' | grep -v 'grpc/gen')

### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest

### Testing
unit-test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

test-race:
	go test ./... --race

### Formatting the code
fmt:
	gofumpt -l -w .
	go mod tidy

check:
	golangci-lint run --timeout=20m0s

### pre commit
pre-commit: fmt check unit-test
	@echo ready to commit...

.PHONY: build