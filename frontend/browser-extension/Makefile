.PHONY: install build clean dev dev-react test dist lint lint-fix icons

ICONS_DIR := src/host/assets/icons
ICON_SVG := $(ICONS_DIR)/icon.svg
ICON_SIZES := 16 32 48 128
ICON_PNGS := $(foreach size,$(ICON_SIZES),$(ICONS_DIR)/icon$(size).png)

install:
	npm install

dist: icons
	npm run dist

build: icons
	npm run build

clean:
	npm run clean
	rm -rf debug-build
	rm -f $(ICON_PNGS)

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

icons: $(ICON_PNGS)

$(ICONS_DIR)/icon%.png: $(ICON_SVG)
	convert -background none -resize $*x$* $< $@
