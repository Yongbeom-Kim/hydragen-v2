setup:
	cp ./.git-hooks/pre-commit ./.git/hooks/pre-commit
	chmod +x ./.git/hooks/pre-commit

start:
	docker compose -f docker-compose.base.yaml -f docker-compose.prod.yaml up

dev:
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml up

dev-attach-flyway:
	docker container exec --env-file .env -it hydragen-v2-postgres-migrate-1 /bin/bash

stop:
	docker compose -f docker-compose.base.yaml -f docker-compose.dev.yaml 
	
clean-volumnes:
	sudo rm -rf ../.hydragen-v2/volumes/*