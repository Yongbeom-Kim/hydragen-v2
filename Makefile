setup:
	cp ./.git-hooks/pre-commit ./.git/hooks/pre-commit
	chmod +x ./.git/hooks/pre-commit

start:
	docker compose -f docker-compose.base.yaml -f docker-compose.prod.yaml up

dev:
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml build
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml up

attach-flyway:
	docker container exec --env-file .env -it hydragen-v2-postgres-migrate-1 /bin/bash

attach-dataset-loader:
	docker container exec --env-file .env -it hydragen-v2-dataset-loader-1 /bin/sh

attach-postgres:
	docker container exec --env-file .env -it hydragen-v2-postgres-1 /bin/bash
	
stop:
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml 
	
clean-volumnes:
	sudo rm -rf ../.hydragen-v2/volumes/*