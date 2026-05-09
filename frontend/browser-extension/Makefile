.PHONY: install build clean dev test dist lint lint-fix icons openapi-client

OPENAPI_URL ?= http://localhost:8080/openapi/web

install:
	npm install

dist: icons
	npm run dist

build: icons
	npm run build

clean:
	npm run clean
	$(MAKE) -C ../shared clean-icons

# Development mode for the Browser Extension (Webpack watch)
# Loads the extension in the browser, watching for file changes.
dev: icons
	npm run dev

openapi-client:
	OPENAPI_URL="$(OPENAPI_URL)" npm run openapi-client

test:
	npm run test

lint:
	npm run lint

check: test lint

lint-fix:
	npm run lint -- --fix

fix: lint-fix

icons:
	$(MAKE) -C ../shared icons
