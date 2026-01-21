PNGS := screenshot.png small_promo_tile.png marquee_promo_tile.png

.PHONY: png

all: png

png: $(PNGS)

screenshot.png: screenshot.md
	mkdir -p dest
	marp $< --images png --allow-local-files -o "dest/$@"

small_promo_tile.png: small_promo_tile.md
	mkdir -p dest
	marp $< --images png --allow-local-files -o "dest/$@"

marquee_promo_tile.png: marquee_promo_tile.md
	mkdir -p dest
	marp $< --images png --allow-local-files -o "dest/$@"

clean:
	rm -rf dest
