#!/bin/sh
set -e
apt-get update
apt-get install -y sudo build-essential wget xxd

cd AppImageKit

VERSION=$(cat .git/ref)-$(uname -m)

bash -ex install-build-deps.sh
bash -ex build.sh

mv build/out AppImageKit-$VERSION
tar czf ../bin/AppImageKit-$VERSION.tgz AppImageKit-$VERSION
