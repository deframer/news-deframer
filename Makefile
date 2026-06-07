DOCKER_REPO := ghcr.io/deframer/news-deframer
BUILD_DIR := bin
CMD_DIR := cmd
DB_IMAGE := pgduckdb/pgduckdb:18-main

ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: all build clean test help coverage lint tidy gen example format-check
.PHONY: infra-env-start infra-env-stop infra-env-down infra-env-zap
.PHONY: docker-all add-feeds service worker thinker thinker-fixer

all: build

infra-env-start:
	$(MAKE) -C infra-env start

infra-env-stop:
	$(MAKE) -C infra-env stop

infra-env-down:
	$(MAKE) -C infra-env down

infra-env-zap:
	$(MAKE) -C infra-env zap

build: gen tidy
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/ ./$(CMD_DIR)/...

clean:
	rm -rf $(BUILD_DIR) gen

test:
	go clean -testcache
	go test -v -race ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: format-check
	golangci-lint run
	gosec -conf .gosec.json ./...
	govulncheck ./...

format-check:
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted:"; \
		gofmt -d .; \
		exit 1; \
	fi

tools-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

goa-install:
	go install goa.design/goa/v3/cmd/goa@latest

check: test lint

tidy: gen
	go mod tidy

gen:
	goa gen github.com/deframer/news-deframer/pkg/design

example:
	goa example github.com/deframer/news-deframer/pkg/design && mkdir -p pkg/service && for f in *.go; do [ -e "$$f" ] || continue; if [ -e "pkg/service/$$f" ]; then rm -f "$$f"; else mv -f "$$f" "pkg/service/$$f"; fi; done

docker-all: $(addprefix docker-,$(notdir $(wildcard build/package/*)))

docker-%:
	@echo "Building Docker image for $*..."
	docker build -t $(DOCKER_REPO)/$*:latest -f build/package/$*/Dockerfile .

import-feeds: build
	./bin/admin feed import -f feeds.json

service-cached: build
	./bin/service

service: build
	./bin/service

migration: build
	./bin/migration

worker: build
	./bin/admin feed sync-all
	./bin/worker

thinker: build
	./bin/worker --mode thinker

thinker-fixer: build
	./bin/worker --mode thinker-fixer

SQL_DIR := sql

$(SQL_DIR)/%.sql: FORCE
	docker run --rm -i --network host \
		-v "$(CURDIR):/workspace" \
		-w /workspace \
		$(DB_IMAGE) \
		psql "postgres://$${DB_USER}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}" \
		-f $@

FORCE:
