.PHONY: lint test fmt bench update-bagend

bench:
	@python3 scripts/bench.py lex

fmt:
	@go fmt ./...

install:
	@go install ./harp/harp

lint:
	@golangci-lint run ./...\
	  --max-issues-per-linter 0\
	  --max-same-issues 0\
	  --allow-parallel-runners

test:
	@go test ./...

update-bagend:
	@go get github.com/mooss/bagend && go mod tidy
