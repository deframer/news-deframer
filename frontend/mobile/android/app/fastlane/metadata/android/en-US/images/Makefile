OUTPUT := featureGraphic.png

.PHONY: feature-graphic clean

feature-graphic: $(OUTPUT)

$(OUTPUT): feature_graphic.md icon.png
	mkdir -p dest
	marp $< --image-scale 1.0 --images png --allow-local-files -o "dest/tmp.png"
	convert "dest/tmp.001.png" -gravity center -crop 1024x500+0+0 +repage "$@"
	rm -f "dest/tmp.001.png"

clean:
	rm -rf dest "$(OUTPUT)"
