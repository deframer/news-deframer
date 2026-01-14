.PHONY: build dist clean test install

install:
	$(MAKE) -C library install

build:
	$(MAKE) -C library build

dist:
	$(MAKE) -C library dist
	cp library/dist/news-deframer-lib.js .

clean:
	$(MAKE) -C library clean

test:
	$(MAKE) -C library test