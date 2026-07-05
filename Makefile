EXAMPLE_DIRS := $(sort $(patsubst %/Makefile,%,$(shell find nilaway -mindepth 2 -name Makefile)))

.PHONY: test test-update list-examples

test: test-update
	@echo "==> git diff"
	pwd
	ls -al
	git diff --exit-code

test-update:
	@set -eux; \
	for dir in $(EXAMPLE_DIRS); do \
		echo "==> $$dir: ci-test"; \
		$(MAKE) -C "$$dir" ci-test; \
	done

list-examples:
	@printf '%s\n' $(EXAMPLE_DIRS)
