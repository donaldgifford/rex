# Rex - Documentation and API Targets
# API documentation generation (OpenAPI/Swagger, Postman) and task management

#######################
##@ Documentation & API

.PHONY: task task-update adr adr-update rfc rfc-update

## Task Management

task: ## Create a new task file interactively
	@ $(MAKE) --no-print-directory log-$@
	@./tools/scripts/tasks/task-create.sh

task-update: ## Update tasks README with current status table
	@ $(MAKE) --no-print-directory log-$@
	@./tools/scripts/tasks/task-update.sh

## Architecture Decision Records (ADR)

adr: ## Create new ADR (usage: make adr "Your ADR Title")
	@ $(MAKE) --no-print-directory log-$@
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "❌ Error: ADR title is required"; \
		echo "Usage: make adr \"Your ADR Title\""; \
		exit 1; \
	fi
	@./tools/docs/create-adr.sh "$(filter-out $@,$(MAKECMDGOALS))"

adr-update: ## Update docs/adr/README.md table
	@ $(MAKE) --no-print-directory log-$@
	@./tools/docs/update-adr-readme.sh

## Request for Comments (RFC)

rfc: ## Create new RFC (usage: make rfc "Your RFC Title")
	@ $(MAKE) --no-print-directory log-$@
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "❌ Error: RFC title is required"; \
		echo "Usage: make rfc \"Your RFC Title\""; \
		exit 1; \
	fi
	@./tools/docs/create-rfc.sh "$(filter-out $@,$(MAKECMDGOALS))"

rfc-update: ## Update docs/rfc/README.md table
	@ $(MAKE) --no-print-directory log-$@
	@./tools/docs/update-rfc-readme.sh
