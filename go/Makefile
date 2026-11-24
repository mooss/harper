.PHONY: lint test

lint:
	@golangci-lint run ./...\
	  --max-issues-per-linter 0\
	  --max-same-issues 0\
	  --allow-parallel-runners

test:
	@go test ./...

fmt:
	@go fmt ./...
