# x-word Go service (gRPC)
BINARY_NAME ?= xword
MAIN_PKG    := ./cmd/xword
GO          := go
GOFLAGS     :=
LDFLAGS     := -s -w
PROTO_ROOT  ?= ../x-proto

.PHONY: all build build-linux test run clean deps lint generate

all: deps build

# Install/update dependencies (creates go.sum)
deps:
	$(GO) mod tidy
	$(GO) mod download

# Build binary (current OS)
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) $(MAIN_PKG)

# Build for Linux (e.g. for containers)
build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PKG)

# Run tests
test:
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage:
	$(GO) test $(GOFLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run the server locally (HTTP_PORT=8080 by default)
run: build
	./bin/$(BINARY_NAME)

# Lint (requires golangci-lint if installed)
lint:
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || $(GO) vet ./...

# Remove build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate Go code from proto (requires protoc, protoc-gen-go, protoc-gen-go-grpc, or run via Docker).
# From repo root: docker run --rm -v $(pwd):/workspace -w /workspace/x-word golang:1.22-bookworm bash -c 'apt-get update -qq && apt-get install -qq -y protobuf-compiler && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0 && protoc -I ../x-proto --go_out=. --go_opt=module=github.com/hungp29/x-word --go-grpc_out=. --go-grpc_opt=module=github.com/hungp29/x-word ../x-proto/word/v1/word.proto'
generate:
	protoc -I $(PROTO_ROOT) --go_out=. --go_opt=module=github.com/hungp29/x-word \
		--go-grpc_out=. --go-grpc_opt=module=github.com/hungp29/x-word \
		$(PROTO_ROOT)/word/v1/word.proto
