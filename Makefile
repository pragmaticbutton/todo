run:
	@go run ./cmd/todo

build-cli:
	@go build -o bin/todo ./cmd/cli

lint:
	@golangci-lint run ./...

test:
	@go test -v -race -cover ./...
