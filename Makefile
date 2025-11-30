run:
	@go run ./cmd/todo

lint:
	@golangci-lint run ./...

test:
	@go test -v -race -cover ./...
