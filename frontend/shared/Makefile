.PHONY: icons clean-icons icons-browser-extension icons-android icons-ios icons-web

ASSETS_DIR := assets
ICON_SVG := $(ASSETS_DIR)/icon.svg

BROWSER_EXTENSION_DIR := $(ASSETS_DIR)/browser-extension
BROWSER_EXTENSION_SIZES := 16 32 48 128
BROWSER_EXTENSION_PNGS := $(foreach size,$(BROWSER_EXTENSION_SIZES),$(BROWSER_EXTENSION_DIR)/icon$(size).png)

ANDROID_DIR := $(ASSETS_DIR)/android
ANDROID_SIZES := 48 72 96 144 192
ANDROID_PNGS := $(foreach size,$(ANDROID_SIZES),$(ANDROID_DIR)/icon-$(size).png)

IOS_DIR := $(ASSETS_DIR)/ios
IOS_PNGS := $(IOS_DIR)/icon-1024.png

WEB_DIR := $(ASSETS_DIR)/web
WEB_PNGS := $(WEB_DIR)/favicon-16.png $(WEB_DIR)/favicon-32.png $(WEB_DIR)/apple-touch-icon.png

icons: icons-browser-extension icons-android icons-ios icons-web

icons-browser-extension: $(BROWSER_EXTENSION_PNGS)

icons-android: $(ANDROID_PNGS)

icons-ios: $(IOS_PNGS)

icons-web: $(WEB_PNGS)

clean-icons:
	rm -f $(BROWSER_EXTENSION_PNGS) $(ANDROID_PNGS) $(IOS_PNGS) $(WEB_PNGS)

$(BROWSER_EXTENSION_DIR)/icon%.png: $(ICON_SVG)
	@mkdir -p $(BROWSER_EXTENSION_DIR)
	rsvg-convert -w $* -h $* $< -o $@

$(ANDROID_DIR)/icon-%.png: $(ICON_SVG)
	@mkdir -p $(ANDROID_DIR)
	rsvg-convert -w $* -h $* $< -o $@

$(IOS_DIR)/icon-1024.png: $(ICON_SVG)
	@mkdir -p $(IOS_DIR)
	rsvg-convert -w 1024 -h 1024 $< -o $@

$(WEB_DIR)/favicon-16.png: $(ICON_SVG)
	@mkdir -p $(WEB_DIR)
	rsvg-convert -w 16 -h 16 $< -o $@

$(WEB_DIR)/favicon-32.png: $(ICON_SVG)
	@mkdir -p $(WEB_DIR)
	rsvg-convert -w 32 -h 32 $< -o $@

$(WEB_DIR)/apple-touch-icon.png: $(ICON_SVG)
	@mkdir -p $(WEB_DIR)
	rsvg-convert -w 180 -h 180 $< -o $@
