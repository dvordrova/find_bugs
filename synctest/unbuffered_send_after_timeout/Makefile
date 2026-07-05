.PHONY: run lint test ci-test

run:
	go run .

lint:
	go test -run TestSynctestFindsBlockedSendAfterTimeout ./...

test:
	go test -run TestOrdinaryTestOnlyChecksTimeout ./...

ci-test: test
	@set -eux; \
	raw="$$(mktemp)"; \
	if go test -run TestSynctestFindsBlockedSendAfterTimeout ./... > "$$raw" 2>&1; then \
		LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
		rm -f "$$raw"; \
		echo "ERROR: synctest check succeeded, expected it to report the blocked send" >&2; \
		exit 1; \
	fi; \
	LC_ALL=C sed -E 's/[[:space:]][0-9.]+s$$/ Xs/g' "$$raw" > lint.logs; \
	rm -f "$$raw"
