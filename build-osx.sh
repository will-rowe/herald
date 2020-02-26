#!/usr/bin/env sh

MAINTAINER=willrowe
APP_NAME=Herald
APP_VERSION=1.0.0
APPDIR=${APP_NAME}.app

mkdir -p $APPDIR/Contents/{MacOS,Resources}
go build -o $APPDIR/Contents/MacOS/$APP_NAME
cat > $APPDIR/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="$APP_VERSION">
<dict>
	<key>CFBundleExecutable</key>
	<string>$APP_NAME</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.$MAINTAINER.$APP_NAME</string>
</dict>
</plist>
EOF
cp gui/assets/icons/icon.icns $APPDIR/Contents/Resources/icon.icns
find $APPDIR
