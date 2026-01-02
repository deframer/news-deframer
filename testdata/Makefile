all:
	@echo all

.PHONY: zap
zap: down start

.PHONY: start
start:
	#ocker compose up -d
	docker compose up -d --build --force-recreate --no-deps

.PHONY: stop
stop:
	docker compose stop

.PHONY: down
down:
	docker compose down --remove-orphans --volumes

.PHONY: wpcli-shell
wpcli-shell:
	docker compose run -it --rm wpcli /bin/bash

.PHONY: clean
clean: down
	docker compose down --rmi local
