APP     := inventory
VERSION ?= $(shell git describe --tags --always)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all build test clean

all: test build

build: ## Build the binary
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/$(APP)

test:
	go vet ./...
	go test -race -count=1 ./...

clean:
	@rm -rf bin/ dist/
	@echo "cleaned $(APP)"
