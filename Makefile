GO ?= go
BIN_DIR ?= bin
CMDS ?= raftkv kvctl

.DEFAULT_GOAL = all

.PHONY: all build test lint vuln clean

all: build lint test

build:
	@mkdir -p $(BIN_DIR)
	@for c in $(CMDS); do \
		echo "go build ./cmd/$$c"; \
		$(GO) build -o $(BIN_DIR)/$$c ./cmd/$$c || exit 1; \
	done

test:
	$(GO) test -race -count=1 ./...

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

clean:
	rm -rf $(BIN_DIR)
