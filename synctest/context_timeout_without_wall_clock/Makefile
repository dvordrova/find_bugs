.PHONY: run lint test ci-test

run:
	go run .

lint:
	go test -run TestSynctestChecksLeaseDeadline ./...

test:
	go test -run TestOrdinaryTestOnlyChecksLeaseStartsActive ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if go test -run TestSynctestChecksLeaseDeadline ./... > "$$raw" 2>&1; then \
		LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: synctest check succeeded, expected it to report the wrong timeout" >&2; \
		exit 1; \
	fi; \
	LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
	rm -f "$$raw"
