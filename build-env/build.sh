#!/bin/sh
set -ex

EXTRA_PKGS="git xxd zsync wget curl libarchive-dev autoconf libtool make gcc libtool libfuse-dev liblzma-dev libglib2.0-dev libssl-dev libinotifytools0-dev liblz4-dev equivs libcairo-dev desktop-file-utils libtool-bin"

apt-get -qq update
apt-get -qq install -y qemu-user-static debootstrap

OUT=./bin/rootfs

debootstrap --foreign --arch=$TARGET_ARCH --variant=buildd $SUITE $OUT
cp /usr/bin/qemu-*-static $OUT/usr/bin/

chroot $OUT /debootstrap/debootstrap --second-stage
chroot $OUT apt-get -qq update
chroot $OUT apt-get -qq install -y $EXTRA_PKGS

# cleanup
rm -f \
    $OUT/var/cache/apt/archives/*.deb \
    $OUT/var/cache/apt/archives/partial/*.deb \
    $OUT/var/cache/apt/*.bin \
    $OUT/var/apt/lists/* || true

cp repo/build-env/* $OUT/../
echo "$SUITE-$TARGET_ARCH" >$OUT/../tag
