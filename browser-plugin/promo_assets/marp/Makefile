SOURCES := $(wildcard *.md)

PNGS  := $(SOURCES:.md=.png)

.PHONY: png

all: png

png:  $(PNGS)

install:
	# echo you also need chrom(ium) for pdf/ppxt
	sudo npm install -g @marp-team/marp-cli
	sudo npm install -g inline-source-cli

%.png: %.md
	mkdir -p dest
	marp $< --images png --allow-local-files --width 1280 --height 800 -o "dest/$@"
	touch "dest/$@"

%: %.md
	$(MAKE) $*.png

clean:
	rm -rf dest
