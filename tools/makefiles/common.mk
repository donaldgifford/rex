# Rex - Common Makefile Definitions
# Shared variables, patterns, and utilities used across all domain makefiles

## Project Variables

PROJECT_NAME  := rex
PROJECT_OWNER := donaldgifford
DESCRIPTION   := Rex - Repo Templating tool with Claude support
PROJECT_URL   := https://github.com/$(PROJECT_OWNER)/$(PROJECT_NAME)
LICENSE       := MIT

## Go Variables

GO          ?= go
GO_PACKAGE  := github.com/$(PROJECT_OWNER)/$(PROJECT_NAME)
GOOS        ?= $(shell $(GO) env GOOS)
GOARCH      ?= $(shell $(GO) env GOARCH)

GOIMPORTS_LOCAL_ARG := -local github.com/donaldgifford

## Build Directories

BIN_DIR        := bin
PLUGIN_BIN_DIR := $(BIN_DIR)/plugins
BUILD_DIR      := $(BIN_DIR)

## Version Information

COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
CUR_VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.0-$(COMMIT_HASH)")

## Build Variables

COVERAGE_OUT := coverage.out


## Mise
MISE ?= mise

###############
##@ Development Setup

.PHONY: setup tools

## Setup and Tools

setup: tools ## Complete development environment setup
	@ $(MAKE) --no-print-directory log-$@
	@echo "✓ Development environment ready"

tools: ## Install Go development tools
	@ $(MAKE) --no-print-directory log-$@
	@ $(MISE) install
	@echo "✓ Go tools installed"

########################################################################
## Self-Documenting Makefile Help                                     ##
## https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html ##
########################################################################

########
##@ Help

.PHONY: help
help:   ## Display this help
	@awk -v "col=\033[36m" -v "nocol=\033[0m" ' \
		BEGIN { FS = ":.*##" ; printf "Usage:\n  make %s<target>%s\n\n", col, nocol } \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  %s%-25s%s %s\n", col, $$1, nocol, $$2 } \
		/^##@/ { printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol } \
	' $(MAKEFILE_LIST)

## Log Pattern
## Automatically logs what a target does by extracting its ## comment
log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN { FS = ":.*?## " }; { printf "\033[36m==> %s\033[0m\n", $$2 }'
