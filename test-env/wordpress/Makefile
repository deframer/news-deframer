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

wpcli:
	docker compose run -it wpcli /bin/bash
