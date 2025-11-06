#  Rex - Go Targets
# Go build, test, lint, and code generation targets

###############
##@ Go Development

.PHONY: build
.PHONY: test test-coverage test-integration
.PHONY: lint lint-fix fmt clean
.PHONY: ci check

## Build Targets

build: build-core ## Build everything


build-core: ## Build core binaries
	@ $(MAKE) --no-print-directory log-$@
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/rex main.go internal/ cmd/
	@echo "✓ Core binaries built"


## Testing

test: ## Run all tests with race detector
	@ $(MAKE) --no-print-directory log-$@
	@go test -v -race ./...

test-coverage: ## Run tests with coverage report
	@ $(MAKE) --no-print-directory log-$@
	@go test -coverprofile=$(COVERAGE_OUT) ./...
	@go tool cover -html=$(COVERAGE_OUT)

test-integration: ## Run all integration tests
	@ $(MAKE) --no-print-directory log-$@
	@for dir in plugins/official/*/; do \
		echo "Testing $$dir..."; \
		go test -v -tags=integration $$dir/... || true; \
	done

## Code Quality

lint: ## Run golangci-lint
	@ $(MAKE) --no-print-directory log-$@
	@golangci-lint run ./...

lint-fix: ## Run golangci-lint with auto-fix
	@ $(MAKE) --no-print-directory log-$@
	@golangci-lint run --fix ./...

fmt: ## Format code with gofmt and goimports
	@ $(MAKE) --no-print-directory log-$@
	@gofmt -s -w .
	@goimports -w $(GOIMPORTS_LOCAL_ARG) .

clean: ## Remove build artifacts
	@ $(MAKE) --no-print-directory log-$@
	@rm -rf $(BIN_DIR)/
	@rm -f $(COVERAGE_OUT)
	@go clean -cache
	@find . -name "*.test" -delete
	@echo "✓ Build artifacts cleaned"


## CI/CD

ci: lint test build ## Run CI pipeline (lint + test + build)
	@ $(MAKE) --no-print-directory log-$@
	@echo "✓ CI pipeline complete"

check: lint test ## Quick pre-commit check (lint + test)
	@ $(MAKE) --no-print-directory log-$@
	@echo "✓ Pre-commit checks passed"
