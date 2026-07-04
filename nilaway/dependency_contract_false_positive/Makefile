# Go tool dependencies. Keep versions pinned; update intentionally.
TOOLS := \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

CUSTOM_LINT := ./custom-gcl

.PHONY: run lint lint-fixed test ci-test tool-update clean

run:
	go run .

lint: $(CUSTOM_LINT)
	$(CUSTOM_LINT) run ./...

lint-fixed: $(CUSTOM_LINT)
	$(CUSTOM_LINT) run --config .golangci.fixed.yaml

test:
	go test ./...

ci-test: tool-update test $(CUSTOM_LINT)
	@set -eu; \
	if $(CUSTOM_LINT) run ./... > lint.logs 2>&1; then \
		echo "ERROR: lint succeeded, expected it to report the example bug" >&2; \
		exit 1; \
	fi; \
	$(CUSTOM_LINT) run --config .golangci.fixed.yaml > lint-fixed.logs 2>&1

# Maintainer: apply pinned tool versions to go.mod/go.sum.
tool-update:
	go get -tool $(TOOLS)
	go mod tidy

# Internal: build local golangci-lint with the NilAway plugin.
$(CUSTOM_LINT): .custom-gcl.yml .golangci.yaml go.mod
	go tool golangci-lint custom

# Maintainer: remove generated custom linter binary.
clean:
	rm -f custom-gcl
