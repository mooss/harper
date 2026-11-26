.PHONY: lint test fmt bench

lint:
	@golangci-lint run ./...\
	  --max-issues-per-linter 0\
	  --max-same-issues 0\
	  --allow-parallel-runners

test:
	@go test ./...

fmt:
	@go fmt ./...

bench:
	@python3 scripts/bench.py lex

update-bagend:
	@go get github.com/mooss/bagend && go mod tidy
