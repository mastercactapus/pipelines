#!/bin/sh
set -e

cd AppImageKit
VERSION=$(cat .git/ref)-$(uname -m)

bash -ex build.sh -n

mv build/out AppImageKit-$VERSION
tar czf ../bin/AppImageKit-$VERSION.tgz AppImageKit-$VERSION
