CUR_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

MINI_PROGRAMS_DIRS := $(shell find $(CUR_DIR) -type f -name go.mod -exec dirname {} \; | xargs -n 1 basename | sort -u)


.PHONY: test-%
test-%:
	@echo "\033[0;32m-- Test $*\033[0m"
	@cd ${CUR_DIR}/$* && make test && echo "\t--- PASS" || echo "\t--- FAILED"

.PHONY: test-all
test-all: $(addprefix test-, $(MINI_PROGRAMS_DIRS))
	@echo "\033[0;32mDone all!\033[0m"

.PHONY: go-mod-tidy-%
go-mod-tidy-%:
	@cd ${CUR_DIR}/$* && make go-mod-tidy

.PHONY: go-mod-tidy-all
go-mod-tidy-all: $(addprefix go-mod-tidy-, $(MINI_PROGRAMS_DIRS))
	@echo "\033[0;32mDone all!\033[0m"
