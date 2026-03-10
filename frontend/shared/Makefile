.PHONY: icons clean-icons

ICONS_DIR := assets/icons
ICON_SVG := $(ICONS_DIR)/icon.svg
ICON_SIZES := 16 32 48 128
ICON_PNGS := $(foreach size,$(ICON_SIZES),$(ICONS_DIR)/icon$(size).png)

icons: $(ICON_PNGS)

clean-icons:
	 rm -f $(ICON_PNGS)

$(ICONS_DIR)/icon%.png: $(ICON_SVG)
	convert -background none -resize $*x$* $< $@
