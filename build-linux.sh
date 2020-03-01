#!/bin/sh

MAINTAINER="willrowe"
APP_NAME=herald
APP_VERSION=1.0.0
APPDIR=${APP}_${APP_VERSION}

mkdir -p $APPDIR/usr/bin
mkdir -p $APPDIR/usr/share/applications
mkdir -p $APPDIR/usr/share/icons/hicolor/1024x1024/apps
mkdir -p $APPDIR/usr/share/icons/hicolor/256x256/apps
mkdir -p $APPDIR/DEBIAN

go build -o $APPDIR/usr/bin/$APP_NAME  -ldflags "-X main.dbLocation=${HOME}/.${APP_NAME}"

cp gui/assets/icons/icon.png $APPDIR/usr/share/icons/hicolor/1024x1024/apps/${APP}.png
cp gui/assets/icons/icon.png $APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png

cat > $APPDIR/usr/share/applications/${APP}.desktop << EOF
[Desktop Entry]
Version=$APP_VERSION
Type=Application
Name=$APP_NAME
Exec=$APP_NAME
Icon=$APP_NAME
Terminal=false
StartupWMClass=$APP_NAME
EOF

cat > $APPDIR/DEBIAN/control << EOF
Package: ${APP_NAME}
Version: $APP_VERSION
Section: base
Priority: optional
Architecture: amd64
Maintainer: $MAINTAINER
Description: -
EOF

dpkg-deb --build $APPDIR
