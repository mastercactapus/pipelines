#!/bin/sh
set -xe

apt-get -qq update
apt-get -qq install -y libc6-dev libssl-dev libexpat1-dev libavcodec-dev libgl1-mesa-dev libqt4-dev

VERSION=$(cat makemkv/VERSION)
OUTDIR=$(pwd)/bin
mkdir -p appdir

tar xf makemkv/makemkv-bin-$VERSION.tar.gz
cd makemkv-bin-$VERSION
mkdir -p tmp
touch tmp/eula_accepted
make
make install DESTDIR=../appdir
cd ..

tar xf makemkv/makemkv-oss-$VERSION.tar.gz
cd makemkv-oss-$VERSION
./configure
make -j$(nproc) >build.log
make install DESTDIR=../appdir
cd ..

# fix version
chmod +x linuxdeployqt/*.AppImage
./linuxdeployqt/*.AppImage --appimage-extract

./squashfs-root/AppRun appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs
./squashfs-root/AppRun appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs -executable=appdir/usr/bin/mmdtsdec
./squashfs-root/AppRun appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs -executable=appdir/usr/bin/makemkvcon
ARCH=x86_64 ./squashfs-root/AppRun appdir/usr/share/applications/makemkv.desktop -appimage

cp MakeMKV-x86_64.AppImage bin/MakeMKV-$VERSION-x86_64.AppImage
