lint:
	golangci-lint run ./...
test:
	go test ./internal/... ./pkg/...