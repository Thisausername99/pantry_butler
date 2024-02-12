SHELL := /bin/zsh

MIGRATION_MOUNT=postgres/migrations
DB_CONNECT_STRING=postgres://postgres:postgres@localhost:5432/pantry_butler_dev?sslmode=disable

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
.DEFAULT_GOAL := help


help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

db_setup: 
	make init_db 
	sleep 2
	make seed_db 
.PHONY: db_setup

init_db: ## Start microservices
	docker-compose up --remove-orphans -d
.PHONY: init_db

seed_db: ## Seed database with mock entries
	migrate -source file://${MIGRATION_MOUNT} -database "${DB_CONNECT_STRING}" up
.PHONY: seed_db


down: ## Down microservices
	migrate -source file://${MIGRATION_MOUNT} -database "${DB_CONNECT_STRING}" drop -f
	docker-compose down --remove-orphans
.PHONY: down
