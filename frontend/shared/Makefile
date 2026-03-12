.PHONY: icons clean-icons icons-browser-extension icons-android icons-ios icons-web sync-mobile-android-icons sync-mobile-ios-icons icons-mobile

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
MOBILE_IOS_APPICON_DIR := ../mobile/ios/NewsDeframerMobile/Images.xcassets/AppIcon.appiconset

WEB_DIR := $(ASSETS_DIR)/web
WEB_PNGS := $(WEB_DIR)/favicon-16.png $(WEB_DIR)/favicon-32.png $(WEB_DIR)/apple-touch-icon.png

MOBILE_ANDROID_RES_DIR := ../mobile/android/app/src/main/res

icons: icons-browser-extension icons-android icons-ios icons-web

icons-browser-extension: $(BROWSER_EXTENSION_PNGS)

icons-android: $(ANDROID_PNGS)

icons-ios: $(IOS_PNGS)

icons-web: $(WEB_PNGS)

sync-mobile-android-icons: icons-android
	install -Dm644 $(ANDROID_DIR)/icon-48.png $(MOBILE_ANDROID_RES_DIR)/mipmap-mdpi/ic_launcher.png
	install -Dm644 $(ANDROID_DIR)/icon-48.png $(MOBILE_ANDROID_RES_DIR)/mipmap-mdpi/ic_launcher_round.png
	install -Dm644 $(ANDROID_DIR)/icon-72.png $(MOBILE_ANDROID_RES_DIR)/mipmap-hdpi/ic_launcher.png
	install -Dm644 $(ANDROID_DIR)/icon-72.png $(MOBILE_ANDROID_RES_DIR)/mipmap-hdpi/ic_launcher_round.png
	install -Dm644 $(ANDROID_DIR)/icon-96.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xhdpi/ic_launcher.png
	install -Dm644 $(ANDROID_DIR)/icon-96.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xhdpi/ic_launcher_round.png
	install -Dm644 $(ANDROID_DIR)/icon-144.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xxhdpi/ic_launcher.png
	install -Dm644 $(ANDROID_DIR)/icon-144.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xxhdpi/ic_launcher_round.png
	install -Dm644 $(ANDROID_DIR)/icon-192.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xxxhdpi/ic_launcher.png
	install -Dm644 $(ANDROID_DIR)/icon-192.png $(MOBILE_ANDROID_RES_DIR)/mipmap-xxxhdpi/ic_launcher_round.png

sync-mobile-ios-icons: icons-ios
	@if ! command -v rsvg-convert >/dev/null 2>&1; then echo "rsvg-convert is required. Install it with: brew install librsvg"; exit 1; fi
	install -d $(MOBILE_IOS_APPICON_DIR)
	rsvg-convert -w 40 -h 40 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-20@2x.png
	rsvg-convert -w 60 -h 60 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-20@3x.png
	rsvg-convert -w 58 -h 58 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-29@2x.png
	rsvg-convert -w 87 -h 87 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-29@3x.png
	rsvg-convert -w 80 -h 80 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-40@2x.png
	rsvg-convert -w 120 -h 120 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-40@3x.png
	rsvg-convert -w 120 -h 120 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-60@2x.png
	rsvg-convert -w 180 -h 180 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-60@3x.png
	rsvg-convert -w 1024 -h 1024 $(ICON_SVG) -o $(MOBILE_IOS_APPICON_DIR)/icon-1024.png
	python3 -c 'import json; from pathlib import Path; path = Path("$(MOBILE_IOS_APPICON_DIR)") / "Contents.json"; data = {"images": [{"filename": "icon-20@2x.png", "idiom": "iphone", "scale": "2x", "size": "20x20"}, {"filename": "icon-20@3x.png", "idiom": "iphone", "scale": "3x", "size": "20x20"}, {"filename": "icon-29@2x.png", "idiom": "iphone", "scale": "2x", "size": "29x29"}, {"filename": "icon-29@3x.png", "idiom": "iphone", "scale": "3x", "size": "29x29"}, {"filename": "icon-40@2x.png", "idiom": "iphone", "scale": "2x", "size": "40x40"}, {"filename": "icon-40@3x.png", "idiom": "iphone", "scale": "3x", "size": "40x40"}, {"filename": "icon-60@2x.png", "idiom": "iphone", "scale": "2x", "size": "60x60"}, {"filename": "icon-60@3x.png", "idiom": "iphone", "scale": "3x", "size": "60x60"}, {"filename": "icon-1024.png", "idiom": "ios-marketing", "scale": "1x", "size": "1024x1024"}], "info": {"author": "xcode", "version": 1}}; path.write_text(json.dumps(data, indent=2) + "\n")'

icons-mobile: sync-mobile-android-icons sync-mobile-ios-icons

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
