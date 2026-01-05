ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: all test-env-start test-env-stop test-env-zap zap start stop down clean

all:
	@echo all

test-env-start:
	$(MAKE) -C test-env start

test-env-stop:
	$(MAKE) -C test-env stop

test-env-zap:
	$(MAKE) -C test-env zap

zap: down start

start:
	#docker compose up -d
	docker compose up -d --build --force-recreate --no-deps

stop:
	docker compose stop

down:
	docker compose down --remove-orphans --volumes

clean: down
	docker compose down --rmi local
	rm -f news-deframer