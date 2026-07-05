.PHONY: run lint test ci-test

run:
	go run .

lint:
	go test -run TestMessageOrderPermutations ./...

test:
	go test -run TestOrdinaryProcessesBrokerOrder ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if go test -count=1 -run TestMessageOrderPermutations ./... > "$$raw" 2>&1; then \
		LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: order permutation check succeeded, expected it to report the order assumption" >&2; \
		exit 1; \
	fi; \
	LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
	rm -f "$$raw"
