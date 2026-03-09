PNGS := screenshot.png small_promo_tile.png marquee_promo_tile.png

.PHONY: png icon

all: png icon

png: $(PNGS)

install:
	# echo you also need chrom(ium) for pdf/ppxt
	sudo npm install -g @marp-team/marp-cli
	sudo npm install -g inline-source-cli
	sudo apt install -y librsvg2-bin

icon: store_icon.png

store_icon.png: assets/store_icon.svg
	mkdir -p dest
	rsvg-convert -w 128 -h 128 $< -o dest/$@

screenshot.png: screenshot.md
	mkdir -p dest
	marp $< --image-scale 1.11111 --images png --allow-local-files -o "dest/tmp-$@"
	for img in dest/tmp-*.png; do \
		final_name=$$(echo $$img | sed 's/tmp-//'); \
		convert "$$img" -gravity Center -crop 1280x800+0+0 +repage "$$final_name"; \
	done
	rm dest/tmp-*.png

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

video_en.png: video_en.md
	mkdir -p dest
	marp $< --image-scale 1.5 --images png --allow-local-files -o "dest/$@"

video_de.png: video_de.md
	mkdir -p dest
	marp $< --image-scale 1.5 --images png --allow-local-files -o "dest/$@"

clean:
	rm -rf dest
