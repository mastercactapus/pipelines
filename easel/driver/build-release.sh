#!/bin/sh
set -e

VERSION=$(cat easel/VERSION)

echo v$VERSION >bin/RELEASE_NAME

cp bin-*/*.AppImage bin/

cat >bin/RELEASE_BODY <<EOF
# Easel Linux Driver v$VERSION

## Usage

The `armhf` build should work for ARMv6 devices and up (Raspberry Pi 2/Zero or newer).


### amd64

- Node.js v$(cat node-amd64/VERSION)
- AppImageKit v$(cat appimagekit-bin-amd64/version)

### armhf

- Node.js v$(cat node-armhf/VERSION)
- AppImageKit v$(cat appimagekit-bin-armhf/version)

EOF
