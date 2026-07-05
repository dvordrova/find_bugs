EXAMPLE_MAKEFILES := $(shell find . -mindepth 3 -name Makefile -not -path './.git/*' -not -path './.agents/*')
EXAMPLE_DIRS := $(sort $(patsubst ./%,%,$(patsubst %/Makefile,%,$(EXAMPLE_MAKEFILES))))
LINT_CONFIG ?= $(config)
ROOT_LINT_CONFIG := $(if $(strip $(LINT_CONFIG)),$(abspath $(LINT_CONFIG)))

.PHONY: test test-update list-examples

test: test-update
	git diff --exit-code

test-update:
	@set -eux; \
	for dir in $(EXAMPLE_DIRS); do \
		echo "==> $$dir: ci-test"; \
		$(MAKE) -C "$$dir" ci-test LINT_CONFIG="$(ROOT_LINT_CONFIG)"; \
	done

list-examples:
	@printf '%s\n' $(EXAMPLE_DIRS)
