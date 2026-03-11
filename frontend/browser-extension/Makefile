.PHONY: install build clean dev dev-react test dist lint lint-fix icons

install:
	npm install

dist: icons
	npm run dist

build: icons
	npm run build

clean:
	npm run clean
	rm -rf debug-build
	$(MAKE) -C ../shared clean-icons

# Development mode for the Browser Extension (Webpack watch)
# Loads the extension in the browser, watching for file changes.
dev: icons
	npm run dev

# Development mode for React Components (Webpack Dev Server)
# Opens a standalone web server to develop UI components in isolation (hot-reload enabled).
dev-react:
	npm run dev-react

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
