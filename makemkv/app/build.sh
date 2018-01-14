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
sed -i "s/^Version=.*/Version=$VERSION/" appdir/usr/share/applications/makemkv.desktop

install -m 755 linuxdeployqt/*.AppImage /usr/bin/linuxdeployqt

linuxdeployqt appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs
linuxdeployqt appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs -executable appdir/usr/bin/mmdtsdec
linuxdeployqt appdir/usr/share/applications/makemkv.desktop -bundle-non-qt-libs -executable appdir/usr/bin/makemkvcon
linuxdeployqt appdir/usr/share/applications/makemkv.desktop -appimage

cp *.AppImage bin/
