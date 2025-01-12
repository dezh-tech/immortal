PACKAGES=$(shell go list ./... | grep -v 'tests' | grep -v 'grpc/gen')

ifneq (,$(filter $(OS),Windows_NT MINGW64))
RM = del /q
else
RM = rm -rf
endif

### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5
	go install github.com/bufbuild/buf/cmd/buf@v1.47

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

### Proto
proto:
	$(RM) infrastructure/grpc_client/gen
	$(RM) delivery/grpc/gen
	cd infrastructure/grpc_client/buf && buf generate --template buf.gen.yaml ../proto
	cd delivery/grpc/buf && buf generate --template buf.gen.yaml ../proto

### pre commit
pre-commit: fmt check unit-test
	@echo ready to commit...

.PHONY: build
