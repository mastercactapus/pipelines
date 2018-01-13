#!/bin/sh
set -e

cd AppImageKit
ca
VERSION=$(cat .git/ref)-$(uname -m)

bash -ex build.sh

mv build/out AppImageKit-$VERSION
tar czf ../bin/AppImageKit-$VERSION.tgz AppImageKit-$VERSION
