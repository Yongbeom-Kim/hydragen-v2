##@ Utility
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Setup
# Prepare local development environment, install hooks, etc.
setup: ## Set up local git hooks
	cp ./.git-hooks/pre-commit ./.git/hooks/pre-commit
	chmod +x ./.git/hooks/pre-commit

##@ Docker Compose
start: ## Start Docker Compose with production config
	docker compose -f docker-compose.base.yaml -f docker-compose.prod.yaml up

dev: ## Build & start Docker Compose with development config
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml build
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml up

attach-flyway: ## Attach a shell to the Flyway migration container
	docker container exec --env-file .env -it hydragen-v2-postgres-migrate-1 /bin/bash

attach-dataset-loader: ## Attach a shell to the dataset-loader container
	docker container exec --env-file .env -it hydragen-v2-dataset-loader-1 /bin/sh

attach-postgres: ## Attach a shell to the Postgres container
	docker container exec --env-file .env -it hydragen-v2-postgres-1 /bin/bash

stop: ## Stop Docker Compose development stack
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml 

##@ Cleaning

clean-volumnes: ## Delete Docker volumes (local development volumes only!)
	sudo rm -rf ../.hydragen-v2/volumes/*

##@ SSH

ssh_root: ## SSH into server as root user
	set -a && . ./.env && set +a && ssh "root@$${PUBLIC_IPV4}"

ssh_app: ## SSH into server as app user
	set -a && . ./.env && set +a && ssh "app@$${PUBLIC_IPV4}"