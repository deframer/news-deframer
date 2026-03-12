# iOS Install

Install `rsvg-convert` via Homebrew:

```sh
brew install librsvg
```

Verify the binary is available:

```sh
which rsvg-convert
rsvg-convert --version
```

Sync shared mobile icons into the iOS asset catalog:

```sh
cd ../shared
make sync-mobile-ios-icons
```
