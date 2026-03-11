.PHONY: install pods start android ios web android-apk android-apk-install android-apk-release android-bundle-release android-remote android-remote-pair android-remote-usb-bootstrap android-metro-reverse ios-build clean test lint check lint-fix fix

ADB_PORT ?= 5555
DEBUG_APK := android/app/build/outputs/apk/debug/app-debug.apk

install:
	npm install

pods:
	bundle install && bundle exec pod install --project-directory=ios

start:
	npm run start

android:
	npm run android

ios:
	npm run ios

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
	cd android && ./gradlew assembleDebug

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

android-apk-release:
	cd android && ./gradlew assembleRelease

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
	bundle install && bundle exec pod install --project-directory=ios && xcodebuild -workspace ios/NewsDeframerMobile.xcworkspace -scheme NewsDeframerMobile -configuration Debug -sdk iphonesimulator -derivedDataPath ios/build

clean:
	rm -rf node_modules ios/build android/app/build android/build web/dist
