PACKAGES=$(shell go list ./... | grep -v 'tests' | grep -v 'grpc/gen')

### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	# TODO ::: go-migrate
	# TODO ::: sqlboiler
	# TODO ::: psql driver

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

### Building
build:
	go build -o build/immortal cmd/main.go

### pre commit
pre-commit: fmt check unit-test
	@echo ready to commit...

### docker-compose
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

### sqlBoiler
 models-generate:
	sqlboiler psql

.PHONY: build
