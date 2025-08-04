ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

all:
	@echo all

.PHONY: deps
deps:
	go install goa.design/goa/v3/cmd/goa@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install go.uber.org/mock/mockgen@latest

.PHONY: build
build:
	$(MAKE) gen
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o infra/deploy/service ./cmd/service
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o infra/deploy/service-cli ./cmd/service-cli
	cp ./gen/http/openapi3.json infra/deploy

.PHONY: run
run: build
	infra/deploy/service

.PHONY: test-ci
test-ci: run-migration test

.PHONY: test
test:
	$(MAKE) gen
	go clean -testcache
	go test -v ./...

.PHONY: lint
lint:
	$(MAKE) gen
	golangci-lint run ./...
	gosec -exclude-dir=gen -exclude-dir=.gomodcache ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: gen
gen:
	#goa gen github.com/egandro/news-deframer/pkg/design
	go mod tidy
	go generate ./...

.PHONY: example
example: gen
	goa example github.com/egandro/news-deframer/pkg/design

.PHONY: package
package:
	docker pull $(CI_REGISTRY_IMAGE):latest || true
	docker pull $(IMAGE)
	cd infra/deploy && \
	DOCKER_BUILDKIT=1 docker build --build-arg "IMAGE=$(IMAGE)"  --cache-from $(CONTAINER_IMAGE):latest \
		--tag $(CI_REGISTRY_IMAGE):$(TAG) \
		--tag $(CI_REGISTRY_IMAGE):latest .
	@echo created $(CI_REGISTRY_IMAGE):$(TAG)

.PHONY: push
push:
	docker push $(CI_REGISTRY_IMAGE):$(TAG)
	docker push $(CI_REGISTRY_IMAGE):latest

.PHONY: run-container
run-container: package
	docker run -p 8080:8080 --rm -it --name $(CI_PROJECT_NAME) $(CI_REGISTRY_IMAGE):latest

