APP_NAME := news-deframer
DOCKER_REPO := egandro
BUILD_DIR := bin
CMD_DIR := cmd
APPS := $(wildcard $(CMD_DIR)/*)
BINS := $(patsubst $(CMD_DIR)/%,$(BUILD_DIR)/%,$(APPS))

.PHONY: all build clean test help

ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: all test-env-start test-env-stop test-env-zap zap build clean test help docker-all docker-build
#start stop down clean

all: build

test-env-start:
	$(MAKE) -C test-env start

test-env-stop:
	$(MAKE) -C test-env stop

test-env-down:
	$(MAKE) -C test-env down

test-env-zap:
	$(MAKE) -C test-env zap

build: $(BINS)

$(BUILD_DIR)/%: $(CMD_DIR)/%
	mkdir -p $(BUILD_DIR)
	go build -o $@ ./$<

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...
	gosec ./...
	govulncheck ./...

tidy:
	go mod tidy

docker-all: $(addprefix docker-,$(notdir $(APPS)))

docker-%:
	@echo "Building Docker image for $*..."
	docker build -t $(DOCKER_REPO)/$*:latest -f build/package/$*/Dockerfile .
