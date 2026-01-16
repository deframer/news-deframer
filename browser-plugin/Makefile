.PHONY: build dist clean test install ext

install:
	$(MAKE) -C library install

build:
	$(MAKE) -C library build

# Build library and copy to root (for Tampermonkey)
dist:
	$(MAKE) -C library dist
	cp library/dist/news-deframer-lib.js .

# Build extension
ext: dist
	cp library/dist/news-deframer-lib.js extension/
	@echo "✅ Extension ready in 'extension/' directory."

# Create a signed CRX file (Drag-and-Drop install)
crx: ext
	@echo "🔑 Checking for private key..."
	@if [ ! -f extension.pem ]; then openssl genrsa -out extension.pem 2048; echo "✅ Generated new key: extension.pem"; fi
	@echo "📦 Packing extension.crx..."
	@npx --yes crx pack extension -o extension.crx -p extension.pem
	@echo "✅ extension.crx is ready. Download and drag into chrome://extensions"

# Create a zip for download (Windows/Remote dev support)
# We zip the FOLDER 'extension' so it extracts cleanly into a subdirectory.
pack: ext
	rm -f extension.zip
	zip -r extension.zip extension -x "extension/.gitignore"
	@echo "📦 extension.zip is ready. Right-click and 'Download' it from VS Code!"

clean:
	$(MAKE) -C library clean
	rm -f news-deframer-lib.js
	rm -f extension/news-deframer-lib.js
	rm -f extension.zip

test:
	$(MAKE) -C library test
