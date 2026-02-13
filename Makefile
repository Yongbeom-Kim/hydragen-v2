##@ Utility
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Setup
# Prepare local development environment, install hooks, etc.
setup: ## Set up local git hooks
	cp ./.git-hooks/pre-commit ./.git/hooks/pre-commit
	chmod +x ./.git/hooks/pre-commit

##@ Docker Compose
COMPOSE_BASE_PROD = -f docker-compose.base.yaml -f docker-compose.prod.yaml
COMPOSE_CX23 = $(COMPOSE_BASE_PROD) -f docker-compose.cx23.yaml
COMPOSE_DEV = -f docker-compose.base.yaml -f docker-compose.dev.yaml

start: ## Start Docker Compose with production config
	docker compose $(COMPOSE_BASE_PROD) build
	docker compose $(COMPOSE_BASE_PROD) up

start_detached: ## Start Docker Compose with production config (detached)
	docker compose $(COMPOSE_BASE_PROD) build
	docker compose $(COMPOSE_BASE_PROD) up -d

start-cx23: ## Start production with cx23 resource limits (2 vCPU, 4GB)
	docker compose $(COMPOSE_CX23) build
	docker compose $(COMPOSE_CX23) up

start_detached-cx23: ## Start production with cx23 resource limits (detached)
	docker compose $(COMPOSE_CX23) build
	docker compose $(COMPOSE_CX23) up -d

dev: ## Build & start Docker Compose with development config
	docker compose $(COMPOSE_DEV) build
	docker compose $(COMPOSE_DEV) up

logs: ## Attach to Docker Compose logs
	docker compose $(COMPOSE_BASE_PROD) logs -f

logs-cx23: ## Attach to Docker Compose logs (cx23 stack)
	docker compose $(COMPOSE_CX23) logs -f

attach-flyway: ## Attach a shell to the Flyway debug container
	docker container exec --env-file .env -it hydragen-v2-postgres-migrate-shell-1 /bin/sh

attach-dataset-loader: ## Attach a shell to the dataset-loader container
	docker container exec --env-file .env -it hydragen-v2-dataset-loader-1 /bin/sh

attach-postgres: ## Attach a shell to the Postgres container
	docker container exec --env-file .env -it hydragen-v2-postgres-1 /bin/bash

stop: ## Stop Docker Compose development and production stacks
	docker compose $(COMPOSE_DEV) down || true
	docker compose $(COMPOSE_BASE_PROD) down || true

##@ Cleaning

clean-volumes: ## Delete Docker volumes (local development volumes only!)
	sudo rm -rf ../.hydragen-v2/volumes/*

##@ SSH

ssh_root: ## SSH into server as root user
	set -a && . ./.env && set +a && ssh "root@$${PUBLIC_IPV4}"

ssh_app: ## SSH into server as app user
	set -a && . ./.env && set +a && ssh "app@$${PUBLIC_IPV4}"
