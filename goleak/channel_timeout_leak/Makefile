# Test-time dependencies. Keep versions pinned; update intentionally.
TEST_DEPS := \
	go.uber.org/goleak@v1.3.0

.PHONY: run lint test ci-test tool-update

run:
	go run .

lint:
	go test -trimpath -tags leakcheck ./...

test:
	go test ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if go test -trimpath -tags leakcheck ./... > "$$raw" 2>&1; then \
		sed -E 's/Goroutine [0-9]+/Goroutine N/g; s/goroutine [0-9]+/goroutine N/g; s/ \+0x[[:xdigit:]]+/ +0xADDR/g; s/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: leak check succeeded, expected it to report the example bug" >&2; \
		exit 1; \
	fi; \
	sed -E 's/Goroutine [0-9]+/Goroutine N/g; s/goroutine [0-9]+/goroutine N/g; s/ \+0x[[:xdigit:]]+/ +0xADDR/g; s/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
	rm -f "$$raw"

# Maintainer: apply pinned test dependency versions to go.mod/go.sum.
tool-update:
	go get $(TEST_DEPS)
	go mod tidy
