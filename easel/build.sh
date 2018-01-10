#!/bin/sh
set -ex

apt-get update
apt-get install -y build-essential python p7zip-full cpio squashfs-tools

EASEL_VERSION=$(cat src/VERSION)
NODE_VERSION=$(cat node/VERSION)


mkdir -p $HOME/.node-gyp/$NODE_VERSION
echo 9 >$HOME/.node-gyp/$NODE_VERSION/installVersion

tar xf node/node*-headers.tar.gz --strip-components=1 -C $HOME/.node-gyp/$NODE_VERSION
tar xf node/node*-linux*.tar.gz --strip-components=1 -C /usr/local

cp src/EaselDriver-$EASEL_VERSION.pkg build-repo/

cd build-repo
make build/iris-lib/iris.js build/node VERSION=$EASEL_VERSION
(cd build/iris-lib/node_modules/serialport && /usr/local/lib/node_modules/npm/bin/node-gyp-bin/node-gyp rebuild)
find build -name node_modules -exec rm -rf {}/.bin \;

cp easel.svg easel-driver.desktop AppRun build/

mksquashfs build build.squashfs -root-owned -noappend

FILE=../bin/EaselDriver-$EASEL_VERSION-x86_64.AppImage
cat ../AppImageKit/runtime-x86_64 build.squashfs >$FILE
chmod a+x $FILE
