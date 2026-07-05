# Go tool dependencies. Keep versions pinned; update intentionally.
TOOLS := \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

LINT_CONFIG ?= $(config)
ACTIVE_LINT_CONFIG = $(if $(LINT_CONFIG),$(LINT_CONFIG),.golangci.yaml)

.PHONY: run lint test ci-test tool-update

run:
	go run .

lint:
	go tool golangci-lint run --config "$(ACTIVE_LINT_CONFIG)" ./...

test:
	go test ./...

ci-test: test
	@set -eux; \
	if [ -n "$(LINT_CONFIG)" ]; then \
		if go tool golangci-lint run --config "$(ACTIVE_LINT_CONFIG)" ./... > lint.logs 2>&1; then \
			printf 'OK: lint succeeded with LINT_CONFIG=%s\n' "$(ACTIVE_LINT_CONFIG)" >> lint.logs; \
		fi; \
		exit 0; \
	fi; \
	if go tool golangci-lint run --config .golangci.yaml ./... > lint.logs 2>&1; then \
		echo "ERROR: lint succeeded, expected sqlclosecheck to report the example bug" >&2; \
		exit 1; \
	fi

# Maintainer: apply pinned tool versions to go.mod/go.sum.
tool-update:
	go get -tool $(TOOLS)
	go mod tidy
