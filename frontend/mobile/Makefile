.PHONY: install pods start android ios ios-preflight ios-sim ios-sim-list web android-apk android-apk-install android-apk-release android-bundle-release android-remote android-remote-pair android-remote-usb-bootstrap android-metro-reverse ios-build ios-sim-build clean test lint check lint-fix fix

ADB_PORT ?= 5555
DEBUG_APK := android/app/build/outputs/apk/debug/app-debug.apk
SIMULATOR ?=
IOS_DERIVED_DATA ?= ios/build
IOS_PROJECT := ios/NewsDeframerMobile.xcodeproj

define check_xcode_ready
	@mkdir -p ios/build
	@rm -f ios/build/xcode-preflight.log
	@xcodebuild -list -project "$(IOS_PROJECT)" >/dev/null 2>ios/build/xcode-preflight.log || { \
		if grep -q "runFirstLaunch" ios/build/xcode-preflight.log; then \
			echo "Xcode is not fully initialized for simulator builds."; \
			echo "Run this once with admin privileges, then retry:"; \
			echo "  sudo xcodebuild -runFirstLaunch"; \
			echo; \
			echo "Raw xcodebuild error:"; \
			cat ios/build/xcode-preflight.log; \
		else \
			cat ios/build/xcode-preflight.log; \
		fi; \
		exit 1; \
	}
endef

define resolve_simulator
	resolved_simulator="$(SIMULATOR)"; \
	if [ -z "$$resolved_simulator" ]; then \
		resolved_simulator=$$(xcrun simctl list devices available | perl -ne 'if (/^\s+(iPhone[^\(]+?) \([0-9A-F-]+\) \((?:Shutdown|Booted)\)\s*$$/) { print $$1; exit }'); \
	fi; \
	if [ -z "$$resolved_simulator" ]; then \
		echo "No available iPhone simulator found."; \
		echo "Use 'make ios-sim-list' to inspect devices or pass SIMULATOR=\"...\"."; \
		exit 1; \
	fi; \
	echo "Using simulator: $$resolved_simulator"; \
	$(1)
endef

install:
	npm install

pods:
	bundle install && bundle exec pod install --project-directory=ios

start:
	npm run start

android:
	npm run android

ios:
	$(MAKE) ios-sim

ios-preflight:
	$(call check_xcode_ready)

ios-sim: ios-preflight
	bundle install && bundle exec pod install --project-directory=ios && $(call resolve_simulator,npm run ios:sim -- --simulator="$$resolved_simulator")

ios-sim-list:
	npm run ios:sim:list

web:
	npm run web

dev: web

test:
	npm run test

lint:
	npm run lint

check: test lint

lint-fix:
	npm run lint -- --fix

fix: lint-fix

android-apk:
	cd android && ./gradlew assembleDebug --rerun-tasks

android-apk-install: android-apk
	@if [ -n "$(ADB_DEVICE)" ]; then adb -s "$(ADB_DEVICE)" install -r "$(DEBUG_APK)"; \
	else \
		devices=$$(adb devices | awk 'NR>1 && $$2=="device" {print $$1}'); \
		count=$$(printf '%s\n' "$$devices" | awk 'NF' | wc -l); \
		if [ "$$count" -eq 1 ]; then \
			device=$$(printf '%s\n' "$$devices" | awk 'NF {print; exit}'); \
			adb -s "$$device" install -r "$(DEBUG_APK)"; \
		elif [ "$$count" -eq 0 ]; then \
			echo "No adb device connected. Connect one device or pass ADB_DEVICE=<ip:port>."; \
			exit 1; \
		else \
			echo "Multiple adb devices detected. Pass ADB_DEVICE=<ip:port>."; \
			adb devices; \
			exit 1; \
		fi; \
	fi

android-metro-reverse:
	@if [ -n "$(ADB_DEVICE)" ]; then adb -s "$(ADB_DEVICE)" reverse tcp:8081 tcp:8081 && adb -s "$(ADB_DEVICE)" reverse --list; \
	else \
		devices=$$(adb devices | awk 'NR>1 && $$2=="device" {print $$1}'); \
		count=$$(printf '%s\n' "$$devices" | awk 'NF' | wc -l); \
		if [ "$$count" -eq 1 ]; then \
			device=$$(printf '%s\n' "$$devices" | awk 'NF {print; exit}'); \
			adb -s "$$device" reverse tcp:8081 tcp:8081 && adb -s "$$device" reverse --list; \
		elif [ "$$count" -eq 0 ]; then \
			echo "No adb device connected. Connect one device or pass ADB_DEVICE=<ip:port>."; \
			exit 1; \
		else \
			echo "Multiple adb devices detected. Pass ADB_DEVICE=<ip:port>."; \
			adb devices; \
			exit 1; \
		fi; \
	fi

android-apk-clean:
	cd android && ./gradlew clean

android-apk-release:
	cd android && ./gradlew app:createBundleReleaseJsAndAssets --rerun-tasks && ./gradlew assembleRelease

android-bundle-release:
	cd android && ./gradlew bundleRelease

android-remote:
	@if [ -z "$(ADB_HOST)" ]; then echo "Usage: make android-remote ADB_HOST=<device-ip> [ADB_PORT=5555]"; exit 1; fi
	adb connect $(ADB_HOST):$(ADB_PORT)
	adb devices

android-remote-pair:
	@if [ -z "$(ADB_PAIR)" ] || [ -z "$(ADB_PAIR_CODE)" ] || [ -z "$(ADB_HOST)" ] || [ -z "$(ADB_PORT)" ]; then echo "Usage: make android-remote-pair ADB_PAIR=<ip:pair-port> ADB_PAIR_CODE=<code> ADB_HOST=<device-ip> ADB_PORT=<connect-port>"; exit 1; fi
	adb pair $(ADB_PAIR) $(ADB_PAIR_CODE)
	adb connect $(ADB_HOST):$(ADB_PORT)
	adb devices

android-remote-usb-bootstrap:
	@if [ -z "$(ADB_HOST)" ]; then echo "Usage: make android-remote-usb-bootstrap ADB_HOST=<device-ip> [ADB_PORT=5555]"; exit 1; fi
	adb tcpip $(ADB_PORT)
	adb connect $(ADB_HOST):$(ADB_PORT)
	adb devices

ios-build:
	$(MAKE) ios-sim-build

ios-sim-build: ios-preflight
	bundle install && bundle exec pod install --project-directory=ios && $(call resolve_simulator,xcodebuild -workspace ios/NewsDeframerMobile.xcworkspace -scheme NewsDeframerMobile -configuration Debug -sdk iphonesimulator -destination "platform=iOS Simulator,name=$$resolved_simulator" -derivedDataPath $(IOS_DERIVED_DATA))

clean:
	rm -rf node_modules ios/build android/app/build android/build web/dist
