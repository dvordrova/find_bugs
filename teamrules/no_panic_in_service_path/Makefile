# Go tool dependencies. Keep versions pinned; update intentionally.
TOOLS := \
	github.com/quasilyte/go-ruleguard/cmd/ruleguard@v0.4.5

RULE_DEPS := \
	github.com/quasilyte/go-ruleguard/dsl@v0.3.23

.PHONY: run lint test ci-test tool-update

run:
	go run .

lint:
	go tool ruleguard -rules rules/no_panic.go ./...

test:
	go test ./...

ci-test: test
	@set -eux; \
	go mod download; \
	raw="$$(mktemp)"; \
	if go tool ruleguard -rules rules/no_panic.go ./... > "$$raw" 2>&1; then \
		sed 's#$(CURDIR)/##g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: lint succeeded, expected ruleguard to report the team rule violation" >&2; \
		exit 1; \
	fi; \
	sed 's#$(CURDIR)/##g' "$$raw" > lint.logs; \
	rm -f "$$raw"

# Maintainer: apply pinned ruleguard versions to go.mod/go.sum.
tool-update:
	go get -tool $(TOOLS)
	go get $(RULE_DEPS)
	go mod tidy
