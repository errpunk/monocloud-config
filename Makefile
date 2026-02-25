BIN := monocloud-config

.PHONY: build test run clean

build:
	go build -o $(BIN) .

test:
	go test ./...

run: build
	./$(BIN)

clean:
	rm -f $(BIN)
