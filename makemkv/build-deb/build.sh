#!/bin/sh
set -xe

VERSION=$(cat src/VERSION)
OUTDIR=$(pwd)/bin
MKPKG="checkinstall --install=no --pkgversion=$VERSION"

ARGS="-y --install=no --pkgversion=$VERSION --pkggroup=makemkv --maintainer=mastercactapus@gmail.com --pakdir $OUTDIR"

tar xfv src/makemkv-bin-$VERSION.tar.gz
cd makemkv-bin-$VERSION
mkdir -p tmp
touch tmp/eula_accepted
make

mkdir -p usr/share/MakeMKV
install -m 644 src/eula_en_linux.txt usr/share/MakeMKV/

echo "usr/share/MakeMKV/eula_en_linux.txt">listfile
checkinstall -D $ARGS --pkgname=makemkv-bin --include=listfile --requires=makemkv-oss

cd ..

tar xfv src/makemkv-oss-$VERSION.tar.gz
cd makemkv-oss-$VERSION
./configure
make
checkinstall -D $ARGS --pkgname=makemkv-oss
