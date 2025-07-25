SHELL := /bin/zsh

MIGRATIONS_PATH?=migrations
MONGODB_URL ?= mongodb://root:password@localhost:27017/pantry_butler_dev?authSource=admin

# Makefile for MongoDB migrations
# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
init_setup: ## Start microservices
	docker compose up --build -d
.PHONY: init_setup




down: ## Down microservices
	migrate -source file://${MIGRATION_MOUNT} -database "${DB_CONNECT_STRING}" drop -f
	docker compose down --remove-orphans
.PHONY: down

gqlgen: ## Generate GraphQL config and code using gqlgen
	go run github.com/99designs/gqlgen generate
.PHONY: gqlgen



# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: help migrate-up migrate-down migrate-reset migrate-status migrate-create migrate-force

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

migrate-up: ## Run all pending migrations
	@echo "$(GREEN)Running migrations up...$(NC)"
	migrate -source file://$(MIGRATIONS_PATH) -database "$(MONGODB_URL)" up

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(MONGODB_URL)" down 1

migrate-reset: ## Rollback all migrations
	@echo "$(RED)Rolling back all migrations...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(MONGODB_URL)" down -all; \
	fi

migrate-status: ## Show current migration status
	@echo "$(GREEN)Current migration status:$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(MONGODB_URL)" version

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
ifndef NAME
	@echo "$(RED)Error: NAME is required$(NC)"
	@echo "Usage: make migrate-create NAME=migration_name"
	@exit 1
endif
	@echo "$(GREEN)Creating migration: $(NAME)$(NC)"
	migrate create -ext json -dir $(MIGRATIONS_PATH) -seq $(NAME)

migrate-force: ## Force set migration version (usage: make migrate-force VERSION=1)
ifndef VERSION
	@echo "$(RED)Error: VERSION is required$(NC)"
	@echo "Usage: make migrate-force VERSION=1"
	@exit 1
endif
	@echo "$(YELLOW)Forcing migration version to $(VERSION)$(NC)"
	migrate -path $(MIGRATIONS_PATH) -database "$(MONGODB_URL)" force $(VERSION)

# Development helpers
dev-setup: ## Set up development database with mock data
	@echo "$(GREEN)Setting up development database...$(NC)"
	make migrate-up

dev-reset: ## Reset development database
	@echo "$(YELLOW)Resetting development database...$(NC)"
	make migrate-reset
	make migrate-up

# Production helpers (use with caution)
prod-migrate: ## Run migrations in production (requires MONGODB_URL)
ifndef MONGODB_URL
	@echo "$(RED)Error: MONGODB_URL environment variable is required for production$(NC)"
	@exit 1
endif
	@echo "$(GREEN)Running production migrations...$(NC)"
	@echo "Database: $(MONGODB_URL)"
	@read -p "Continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(MONGODB_URL)" up; \
	fi