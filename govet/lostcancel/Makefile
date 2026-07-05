# Go tool dependencies. Keep versions pinned; update intentionally.
TOOLS := \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

.PHONY: run lint test ci-test tool-update

run:
	go run .

lint:
	go tool golangci-lint run --config .golangci.yaml ./...

test:
	go test ./...

ci-test: test
	@set -eux; \
	if go tool golangci-lint run --config .golangci.yaml ./... > lint.logs 2>&1; then \
		echo "ERROR: lint succeeded, expected govet to report the example bug" >&2; \
		exit 1; \
	fi

# Maintainer: apply pinned tool versions to go.mod/go.sum.
tool-update:
	go get -tool $(TOOLS)
	go mod tidy
