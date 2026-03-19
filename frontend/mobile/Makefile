.PHONY: install pods start android ios ios-preflight ios-sim ios-sim-list web android-apk android-apk-install android-apk-release android-apk-release-install android-bundle-release android-remote-qr android-metro-reverse ios-build ios-sim-build clean test check-mobile-workflows lint check lint-fix fix

DEBUG_APK := android/app/build/outputs/apk/debug/app-debug.apk
RELEASE_APK := android/app/build/outputs/apk/release/app-release.apk
SIMULATOR ?=
IOS_DERIVED_DATA ?= ios/build
ADB_WAIT_TIMEOUT ?= 90

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
	bash ./scripts/ios-preflight.sh

ios-sim: ios-preflight
	bash ./scripts/ios-sim.sh "$(SIMULATOR)"

ios-sim-list:
	npm run ios:sim:list

web:
	npm run web

dev: web

test:
	npm run test

check-mobile-workflows:
	bash ./scripts/check-mobile-workflows.sh

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
		bash ./scripts/adb-single-device.sh install -r "$(DEBUG_APK)"; \
	fi

android-metro-reverse:
	@if [ -n "$(ADB_DEVICE)" ]; then adb -s "$(ADB_DEVICE)" reverse tcp:8081 tcp:8081 && adb -s "$(ADB_DEVICE)" reverse --list; \
	else \
		bash ./scripts/adb-single-device.sh reverse tcp:8081 tcp:8081 && bash ./scripts/adb-single-device.sh reverse --list; \
	fi

android-apk-clean:
	cd android && ./gradlew clean

android-apk-release:
	cd android && ./gradlew app:createBundleReleaseJsAndAssets --rerun-tasks && ./gradlew assembleRelease

android-apk-release-install: android-apk-release
	@if [ -n "$(ADB_DEVICE)" ]; then adb -s "$(ADB_DEVICE)" install -r "$(RELEASE_APK)"; \
	else \
		bash ./scripts/adb-single-device.sh install -r "$(RELEASE_APK)"; \
	fi

android-bundle-release:
	cd android && ./gradlew bundleRelease

android-remote-qr:
	@ADB_CONNECT="$(ADB_CONNECT)" ADB_WAIT_TIMEOUT="$(ADB_WAIT_TIMEOUT)" bash ./scripts/android-remote-qr.sh

ios-build:
	$(MAKE) ios-sim-build

ios-sim-build: ios-preflight
	bash ./scripts/ios-sim-build.sh "$(SIMULATOR)" "$(IOS_DERIVED_DATA)"

clean:
	rm -rf node_modules ios/build android/app/build android/build web/dist
