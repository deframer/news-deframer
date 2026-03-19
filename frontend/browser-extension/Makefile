.PHONY: install build clean dev test dist lint lint-fix icons

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
