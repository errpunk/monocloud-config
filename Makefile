BIN := monocloud-config
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build test run clean

build:
	go build $(LDFLAGS) -o $(BIN) .

test:
	go test ./...

run: build
	./$(BIN)

clean:
	rm -f $(BIN)
