.PHONY: run lint test ci-test

run:
	go run .

lint:
	go test -run TestSynctestChecksBeforeCancel ./...

test:
	go test -run TestOrdinaryTestOnlyChecksAfterCancel ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if go test -run TestSynctestChecksBeforeCancel ./... > "$$raw" 2>&1; then \
		LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: synctest check succeeded, expected it to report the early callback" >&2; \
		exit 1; \
	fi; \
	LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
	rm -f "$$raw"
