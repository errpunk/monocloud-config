BIN := bin/monocloud-config
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
IMAGE := errpunk/monocloud-config

.PHONY: build build-linux docker-build docker-push test run clean

build:
	mkdir -p bin
	go build $(LDFLAGS) -o $(BIN) .

build-linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BIN)-linux-arm64 .

# Build image for the local platform (fast, for testing)
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

# Build and push multi-platform image (linux/amd64 + linux/arm64)
docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		-t $(IMAGE):$(VERSION) -t $(IMAGE):latest \
		--push .

test:
	go test ./...

run: build
	./$(BIN)

clean:
	rm -rf bin/
