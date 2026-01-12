.PHONY: all zap start stop down clean

all:
	@echo all

zap: down start

start:
	docker compose up -d --build --force-recreate --no-deps

stop:
	docker compose stop

down:
	docker compose down --remove-orphans --volumes

clean: down
	docker compose down --rmi local

postgres-bash:
	docker compose exec postgres bash

postgres-logs:
	docker compose logs -f postgres

psql:
	docker compose exec -it postgres psql -U $(DB) -d $(DB)

psql-deframer:
	DB=deframer $(MAKE) psql
