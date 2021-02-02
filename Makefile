L_PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
L_LOCAL_DIR   := $(L_PROJECT_DIR)/.local

.PHONY: all
all:
	@echo "Available commands:"
	@echo ""
	@echo "	codecov		Upload a codecov report"
	@echo ""

.PHONY: codecov
codecov:
	mkdir -p "$(L_LOCAL_DIR)"
	go test -race -coverprofile="$(L_LOCAL_DIR)/coverage.txt" -covermode=atomic
	curl -s https://codecov.io/bash | bash -s - -t $(KINIT_CODECOV_TOKEN) -f "$(L_LOCAL_DIR)/coverage.txt"
