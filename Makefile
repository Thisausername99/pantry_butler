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
	docker compose up --remove-orphans -d
.PHONY: init_setup

down: ## Down microservices
	docker compose down --remove-orphans
.PHONY: down

gqlgen: ## Generate GraphQL config and code using gqlgen
	go run github.com/99designs/gqlgen generate
.PHONY: gqlgen
