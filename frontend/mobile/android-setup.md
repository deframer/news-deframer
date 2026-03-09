# Android

## Windows Setup

### Todo

- Automate some steps e.g. license acceptance
- Automate the generation of a new emulator

### Reinstall notice

You might want to delete/rename these directories.

```cmd
"%LOCALAPPDATA%\Android"
"%USERPROFILE%\.android"
"%USERPROFILE%\.gradle"
"%APPDATA%"\Google\AndroidStudio*
```

### Install Hypervisor

As Admin (you might need to reboot the computer after this)

```powershell
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
```

### Install Chocolatey

<https://chocolatey.org/install>

```powershell
rem you need to be admin
rem close the shell after this
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

### Packages from Chocolatey

```powershell
rem choco install -y googlechrome --ignore-checksums
rem choco install -y vscode
choco install -y 7zip
choco install -y git
choco install -y wget
choco install -y curl
choco install -y winscp
choco install -y openjdk
choco install -y androidstudio
```

### Android SDK (only)

- you can do this via the UI of Android Studio

#### Android SDK (cmdline-tools)

```cmd
mkdir "%LOCALAPPDATA%\Android"
mkdir "%LOCALAPPDATA%\Android\sdk"
set ANDROID_HOME=%LOCALAPPDATA%\Android\sdk
cd "%ANDROID_HOME%"

mkdir cmdline-tools
cd cmdline-tools
set COMMANDLINE_TOOLS_VERSION=13114758_latest
wget -c -t0 https://dl.google.com/android/repository/commandlinetools-win-%COMMANDLINE_TOOLS_VERSION%.zip
7z x commandlinetools-win-%COMMANDLINE_TOOLS_VERSION%.zip
del commandlinetools-win-%COMMANDLINE_TOOLS_VERSION%.zip
ren cmdline-tools latest
set PATH=%ANDROID_HOME%\cmdline-tools\latest\bin;%PATH%;
sdkmanager --update
rem you need to start and accept the licenses
sdkmanager --licenses
sdkmanager --list
```

#### Android SDK (via sdkmanager)

```cmd
rem emulator + android sdk 33
set ANDROID_HOME=%LOCALAPPDATA%\Android\sdk
set PATH=%ANDROID_HOME%\cmdline-tools\latest\bin;%PATH%;
sdkmanager "platform-tools" "platforms;android-33" "build-tools;33.0.3" "emulator" "extras;android;m2repository"
sdkmanager "system-images;android-33;google_apis;x86_64"
```

#### Create Android AVD (via avdmanager)

```cmd
rem emulator + android sdk 33
set ANDROID_HOME=%LOCALAPPDATA%\Android\sdk
set PATH=%ANDROID_HOME%\cmdline-tools\latest\bin;%PATH%;

set AVD_NAME=testemu
rem you might want to add "--force"
avdmanager create avd -n %AVD_NAME% -k "system-images;android-33;google_apis;x86_64" -c 256M -d pixel_4_xl
avdmanager list avd
rem avdmanager delete avd -n %AVD_NAME%
```

#### Start Android AVD (via emulator)

```cmd
rem emulator + android sdk 33
set ANDROID_HOME=%LOCALAPPDATA%\Android\sdk
set PATH=%ANDROID_HOME%\cmdline-tools\latest\bin;%PATH%;
set PATH=%ANDROID_HOME%\emulator;%PATH%;

emulator list
set AVD_NAME=testemu
emulator -avd %AVD_NAME% -no-metrics -memory 2048 -no-audio -no-boot-anim -gpu auto -cores 4
```

#### Reset ADV Snapshot

- Start with no snapshot `emulator -avd %AVD_NAME% -no-snapshot-load`
- Reboot with adb `adb -e reboot`

### Helper (script for android sdk enviornment)

Create a batch file on your desktop: `android-env.cmd`

```txt
set ANDROID_HOME=%LOCALAPPDATA%\Android\sdk
set PATH=%ANDROID_HOME%\cmdline-tools\latest\bin;%PATH%;
set PATH=%ANDROID_HOME%\emulator;%PATH%;
set PATH=%ANDROID_HOME%\tools;%PATH%;
set PATH=%ANDROID_HOME%\platform-tools;%PATH%;

cmd
```

## Linux Setup

```bash
apt-get update
apt-get install -y default-jre default-jdk
apt-get install -y wget unzip

mkdir -p ~/Android/sdk
export ANDROID_HOME=~/Android/sdk
export "PATH=$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"
export "PATH=$ANDROID_HOME/emulator:$PATH"
export "PATH=$ANDROID_HOME/tools:$PATH"
export "PATH=$ANDROID_HOME/platform-tools:$PATH"

cd "${ANDROID_HOME}"
mkdir -p cmdline-tools
cd cmdline-tools
export COMMANDLINE_TOOLS_VERSION=13114758_latest
wget -c -t0 "https://dl.google.com/android/repository/commandlinetools-linux-${COMMANDLINE_TOOLS_VERSION}.zip"
unzip "commandlinetools-linux-${COMMANDLINE_TOOLS_VERSION}.zip"
rm -f "commandlinetools-linux-${COMMANDLINE_TOOLS_VERSION}.zip"
mv cmdline-tools latest
sdkmanager --update
sdkmanager --licenses
sdkmanager --list
```

### Android SDK (via sdkmanager - Linux)

```bash
sdkmanager "platform-tools" "platforms;android-36" "build-tools;36.1.0" "emulator" "extras;android;m2repository"
sdkmanager "system-images;android-36;google_apis;x86_64"
```

#### Create Android AVD (via avdmanager  - Linux)

```bash
export AVD_NAME=testemu
avdmanager create avd -n "${AVD_NAME}" -k "system-images;android-36;google_apis;x86_64" -c 256M -d pixel_4_xl
avdmanager list avd
# you might want to add "--force"
# avdmanager delete avd -n "${AVD_NAME}"
```

## MacOS

```zsh
brew install openjdk
sudo ln -sfn $HOMEBREW_PREFIX/opt/openjdk/libexec/openjdk.jdk /Library/Java/JavaVirtualMachines/openjdk.jdk
```

```zsh
mkdir -p ~/Library/Android/sdk
export ANDROID_HOME=~/Library/Android/sdk
export "PATH=$ANDROID_HOME/cmdline-tools/latest/bin:$PATH"
export "PATH=$ANDROID_HOME/emulator:$PATH"
export "PATH=$ANDROID_HOME/tools:$PATH"
export "PATH=$ANDROID_HOME/platform-tools:$PATH"

cd "${ANDROID_HOME}"
mkdir -p cmdline-tools
cd cmdline-tools
export COMMANDLINE_TOOLS_VERSION=13114758_latest
wget -c -t0 "https://dl.google.com/android/repository/commandlinetools-mac-${COMMANDLINE_TOOLS_VERSION}.zip"
unzip "commandlinetools-mac-${COMMANDLINE_TOOLS_VERSION}.zip"
rm -f "commandlinetools-mac-${COMMANDLINE_TOOLS_VERSION}.zip"
mv cmdline-tools latest
sdkmanager --update
# you need to start and accept the licenses
sdkmanager --licenses
sdkmanager --list
```

### Android SDK (via sdkmanager - Mac)

```bash
sdkmanager "platform-tools" "platforms;android-36" "build-tools;36.1.0" "emulator" "extras;android;m2repository"
sdkmanager "system-images;android-36;google_apis;arm64-v8a"
```

#### Create Android AVD (via avdmanager  - Mac)

```bash
export AVD_NAME=testemu
# you might want to add "--force"
avdmanager create avd -n "${AVD_NAME}" -k "system-images;android-36;google_apis;arm64-v8a" -c 256M -d pixel_4_xl
avdmanager list avd
# avdmanager delete avd -n "${AVD_NAME}"
```
