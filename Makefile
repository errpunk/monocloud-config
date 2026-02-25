BIN := monocloud-config
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build build-linux test run clean

build:
	go build $(LDFLAGS) -o $(BIN) .

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BIN)-linux-arm64 .

test:
	go test ./...

run: build
	./$(BIN)

clean:
	rm -f $(BIN) $(BIN)-linux-amd64 $(BIN)-linux-arm64
