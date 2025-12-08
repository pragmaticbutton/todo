run:
	@go run ./cmd/todo

build-cli:
	@go build -o bin/todo ./cmd/cli

lint:
	@golangci-lint run --config .github/linters/.golangci.yaml ./...

test:
	@go test -v -race -cover ./...

golden:
	go test -v -race -cover ./internal/cli/... -- -update
