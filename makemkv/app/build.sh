#!/bin/sh
set -xe

apt-get -qq update
apt-get -qq install -y libc6-dev libssl-dev libexpat1-dev libavcodec-dev libgl1-mesa-dev libqt4-dev

VERSION=$(cat src/VERSION)
OUTDIR=$(pwd)/bin

tar xfv makemkv/makemkv-bin-$VERSION.tar.gz
cd makemkv-bin-$VERSION
mkdir -p tmp
touch tmp/eula_accepted
make

mkdir -p usr/share/MakeMKV
install -m 644 src/eula_en_linux.txt usr/share/MakeMKV/

tar xfv makemkv/makemkv-oss-$VERSION.tar.gz
cd makemkv-oss-$VERSION
./configure
make
