.PHONY: install build clean dev test dist lint lint-fix icons

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
	rm -f $(ICON_PNGS)

dev: icons
	npm run dev

test:
	npm run test

lint:
	npm run lint

lint-fix:
	npm run lint -- --fix

icons: $(ICON_PNGS)

$(ICONS_DIR)/icon%.png: $(ICON_SVG)
	convert -background none -resize $*x$* $< $@
