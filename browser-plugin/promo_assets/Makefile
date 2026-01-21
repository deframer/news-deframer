PNGS := screenshot.png small_promo_tile.png marquee_promo_tile.png

.PHONY: png

all: png

png: $(PNGS)

screenshot.png: screenshot.md
	mkdir -p dest
	marp $< --images png --allow-local-files -o "dest/$@"

small_promo_tile.png: small_promo_tile.md
	mkdir -p dest
	marp $< --image-scale 0.39 --images png --allow-local-files -o "dest/tmp.png"
	convert "dest/tmp.001.png" -gravity center -extent 440x280 "dest/$@"
	rm -f "dest/tmp.001.png"

marquee_promo_tile.png: marquee_promo_tile.md
	mkdir -p dest
	marp $< --image-scale 1.1 --images png --allow-local-files -o "dest/tmp.png"
	convert "dest/tmp.001.png" -gravity center -crop 1400x560+0+0 +repage "dest/$@"
	rm -f "dest/tmp.001.png"

clean:
	rm -rf dest
