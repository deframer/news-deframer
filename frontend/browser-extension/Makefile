.PHONY: install build clean dev test dist lint lint-fix

install:
	npm install

dist:
	npm run dist

build:
	npm run build

clean:
	npm run clean

# Development mode for the Browser Extension (Webpack watch)
# Loads the extension in the browser, watching for file changes.
dev:
	npm run dev

test:
	npm run test

lint:
	npm run lint

check: test lint

lint-fix:
	npm run lint -- --fix

fix: lint-fix
