.PHONY: run lint test ci-test

run:
	go run .

lint:
	GOMAXPROCS=1 go test -trimpath -race ./...

test:
	GOMAXPROCS=1 go test ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if GOMAXPROCS=1 go test -trimpath -race ./... > "$$raw" 2>&1; then \
		LC_ALL=C sed -E 's/0x[^[:space:]]+/0xADDR/g; s/goroutine [0-9]+/goroutine N/g; s/Goroutine [0-9]+/Goroutine N/g; s/\([0-9.]+s\)/(Xs)/g; s/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > race.logs; \
		rm -f "$$raw"; \
		echo "ERROR: race detector succeeded, expected it to report the example bug" >&2; \
		exit 1; \
	fi; \
	LC_ALL=C sed -E 's/0x[^[:space:]]+/0xADDR/g; s/goroutine [0-9]+/goroutine N/g; s/Goroutine [0-9]+/Goroutine N/g; s/\([0-9.]+s\)/(Xs)/g; s/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > race.logs; \
	rm -f "$$raw"
