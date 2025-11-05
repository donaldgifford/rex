# Rex - Root Makefile
# Orchestration and domain-specific makefile includes
#
# This root Makefile coordinates domain-specific makefiles:
#   - makefiles/common.mk    - Shared variables and patterns
#   - makefiles/go.mk        - Go backend targets
#   - makefiles/ui.mk        - Web UI targets
#   - makefiles/docker.mk    - Docker and LocalStack
#   - makefiles/db.mk        - Database migrations
#   - makefiles/docs.mk      - API docs and task management
#   - makefiles/demo.mk      - E2E demos and MVP testing
#
# All domain targets are defined in their respective makefiles.
# This file only contains orchestration targets that coordinate across domains.

.DEFAULT_GOAL := help

## Include Domain Makefiles
## Order matters: common.mk must be first (defines shared variables and help system)

include tools/makefiles/common.mk
include tools/makefiles/go.mk
include tools/makefiles/docs.mk
include tools/makefiles/demo.mk

######################
##@ Orchestration

.PHONY: all clean-all

all: lint test build build-ui ## Build everything (lint + test + backend + ui)
	@ $(MAKE) --no-print-directory log-$@
	@echo "✓ All targets complete"

clean-all: clean api-docs-clean ## Clean all build artifacts and generated files
	@ $(MAKE) --no-print-directory log-$@
	@echo "✓ All artifacts cleaned"

## Special Pattern to Capture Remaining Arguments
## Allows commands like: make adr "My Title" or make rfc "My RFC"
## This must be at the end of the Makefile
%:
	@:

## End of Root Makefile
## See makefiles/*.mk for domain-specific targets
