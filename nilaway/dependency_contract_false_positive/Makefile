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

ci-test: test $(CUSTOM_LINT)
	@set -eux; \
	export GOLANGCI_LINT_CACHE="$$(mktemp -d)"; \
	export NO_COLOR=1; \
	export TERM=dumb; \
	if $(CUSTOM_LINT) --color=never run ./... > lint.logs 2>&1; then \
		echo "ERROR: lint succeeded, expected it to report the example bug" >&2; \
		exit 1; \
	fi; \
	if ! $(CUSTOM_LINT) --color=never run --config .golangci.fixed.yaml > lint-fixed.logs 2>&1; then \
		cat lint-fixed.logs >&2; \
		exit 1; \
	fi

# Maintainer: apply pinned tool versions to go.mod/go.sum.
tool-update:
	go get -tool $(TOOLS)
	go mod tidy

# Internal: build local golangci-lint with the NilAway plugin.
$(CUSTOM_LINT): .custom-gcl.yml .golangci.yaml go.mod
	tmpcache="$$(mktemp -d)"; \
	export GOLANGCI_LINT_CACHE="$$tmpcache"; \
	go tool golangci-lint custom

# Maintainer: remove generated custom linter binary.
clean:
	rm -f custom-gcl
