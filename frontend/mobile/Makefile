.PHONY: install pods start android ios web android-apk android-apk-release android-bundle-release ios-build clean test lint check lint-fix fix

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

android-apk-release:
	cd android && ./gradlew assembleRelease

android-bundle-release:
	cd android && ./gradlew bundleRelease

ios-build:
	bundle install && bundle exec pod install --project-directory=ios && xcodebuild -workspace ios/NewsDeframerMobile.xcworkspace -scheme NewsDeframerMobile -configuration Debug -sdk iphonesimulator -derivedDataPath ios/build

clean:
	rm -rf node_modules ios/build android/app/build android/build web/dist
