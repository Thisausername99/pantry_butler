SHELL := /bin/zsh
# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color


# Makefile for MongoDB migrations
# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
init_setup: ## Start services
	@if [ "$(BUILD)" = "true" ]; then \
		echo "$(GREEN)Building and starting services...$(NC)"; \
		docker compose up --build --remove-orphans -d; \
	else \
		echo "$(GREEN)Starting services...$(NC)"; \
		docker compose up --remove-orphans -d; \
	fi
.PHONY: init_setup

down: ## Down microservices
	docker compose down --remove-orphans
.PHONY: down

gqlgen: ## Generate GraphQL config and code using gqlgen
	go run github.com/99designs/gqlgen generate
.PHONY: gqlgen


# Generate Gomock for interface
generate_gomock: ## Generate Gomock for interface
	go generate ./internal/mocks
.PHONY: generate_gomock


# Run tests
run_tests: ## Run tests
	go test -v ./internal/usecase/test/...
	go test -v ./internal/adapter/persistence/mongo/test/...
.PHONY: run_tests

# Universal rebuild service (usage: make rebuild SERVICE=server)
rebuild: ## Rebuild and restart any service (usage: make rebuild SERVICE=server)
	docker stop $(SERVICE)
	docker compose build $(SERVICE)
	docker compose up -d $(SERVICE)
.PHONY: rebuild



