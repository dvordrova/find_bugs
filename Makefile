EXAMPLE_MAKEFILES := $(shell find . -mindepth 3 -name Makefile -not -path './.git/*' -not -path './.agents/*')
EXAMPLE_DIRS := $(sort $(patsubst ./%,%,$(patsubst %/Makefile,%,$(EXAMPLE_MAKEFILES))))

.PHONY: test test-update list-examples

test: test-update
	git diff --exit-code

test-update:
	@set -eux; \
	for dir in $(EXAMPLE_DIRS); do \
		echo "==> $$dir: ci-test"; \
		$(MAKE) -C "$$dir" ci-test; \
	done

list-examples:
	@printf '%s\n' $(EXAMPLE_DIRS)
