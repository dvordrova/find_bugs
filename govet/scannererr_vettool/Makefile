# Go analysis dependencies. Keep versions pinned; update intentionally.
ANALYZER_MODULES := \
	golang.org/x/tools@v0.47.0

VETTOOL := ./scannererr-vettool

.PHONY: run lint test ci-test tool-update clean

run:
	go run .

lint: $(VETTOOL)
	go vet -vettool=$(VETTOOL) ./...

test:
	go test ./...

ci-test: test $(VETTOOL)
	@set -eux; \
	if go vet -vettool=$(VETTOOL) ./... > lint.logs 2>&1; then \
		echo "ERROR: lint succeeded, expected scannererr to report the example bug" >&2; \
		exit 1; \
	fi

# Maintainer: apply pinned analyzer module versions to go.mod/go.sum.
tool-update:
	go get $(ANALYZER_MODULES)
	go mod tidy

# Internal: build local vettool with the scannererr analyzer.
$(VETTOOL): tools/scannererr/main.go go.mod
	go build -o $(VETTOOL) ./tools/scannererr

# Maintainer: remove generated vettool binary.
clean:
	rm -f $(VETTOOL)
