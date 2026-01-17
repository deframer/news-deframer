APP_NAME := news-deframer
DOCKER_REPO := egandro
BUILD_DIR := bin
CMD_DIR := cmd
DOCKER_COMPOSE_FILE ?= docker-compose.yml

# we want to avoid to share the developer .env file with docker compose (the DSN hosts etc. are different)
COMPOSE_ENV_FILE ?= .env-compose
DOCKER_ENV_FLAG := $(if $(wildcard $(COMPOSE_ENV_FILE)),--env-file $(COMPOSE_ENV_FILE),--env-file /dev/null)

ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: all build clean test help

.PHONY: all service worker test-env-start test-env-stop test-env-down test-env-zap infra-env-start infra-env-stop infra-env-down infra-env-zap zap build clean test help docker-all docker-build

all: build

test-env-start:
	$(MAKE) -C test-env start

test-env-stop:
	$(MAKE) -C test-env stop

test-env-down:
	$(MAKE) -C test-env down

test-env-zap:
	$(MAKE) -C test-env zap

infra-env-start:
	$(MAKE) -C infra-env start

infra-env-stop:
	$(MAKE) -C infra-env stop

infra-env-down:
	$(MAKE) -C infra-env down

infra-env-zap:
	$(MAKE) -C infra-env zap

zap: down start

# DOCKER_COMPOSE_FILE=docker-compose-lb.yml make start/stop/down/logs

start:
	docker compose $(DOCKER_ENV_FLAG) -f $(DOCKER_COMPOSE_FILE) up -d --build --force-recreate --no-deps

stop:
	docker compose $(DOCKER_ENV_FLAG) -f $(DOCKER_COMPOSE_FILE) stop

down:
	docker compose $(DOCKER_ENV_FLAG) -f $(DOCKER_COMPOSE_FILE) down --remove-orphans --volumes

logs:
	docker compose $(DOCKER_ENV_FLAG) -f $(DOCKER_COMPOSE_FILE) logs -f

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/ ./$(CMD_DIR)/...

clean:
	rm -rf $(BUILD_DIR)
	docker compose $(DOCKER_ENV_FLAG) -f $(DOCKER_COMPOSE_FILE) down --rmi local

test:
	go clean -testcache
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...
	gosec ./...
	govulncheck ./...
	gofmt -l .
	cd browser-plugin && npm run lint

tidy:
	go mod tidy

docker-all: $(addprefix docker-,$(notdir $(wildcard build/package/*)))

docker-%:
	@echo "Building Docker image for $*..."
	docker build -t $(DOCKER_REPO)/$*:latest -f build/package/$*/Dockerfile .

add-feeds: build
	@if [ ! -f feeds.json ]; then echo "feeds.json not found"; exit 1; fi
	@jq -r '.[]' feeds.json | while read url; do \
		echo "Adding feed: $$url"; \
		./bin/admin feed add --enabled --polling "$$url"; \
	done

service: build
	./bin/service

worker: build
	./bin/admin feed sync-all
	./bin/worker
