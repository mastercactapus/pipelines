#!/bin/sh
set -e

cd AppImageKit
VERSION=$(cat .git/ref)-$(dpkg --print-architecture)

echo "Sending build output to build.log"
bash -ex build.sh -n >build.log

set -x

mv out AppImageKit-$VERSION
tar czf ../bin/AppImageKit-$VERSION.tgz AppImageKit-$VERSION
