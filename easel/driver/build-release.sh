#!/bin/sh
set -e

VERSION=$(cat easel/VERSION)

echo v$VERSION >bin/RELEASE_NAME

cp bin-*/*.AppImage bin/

cat >bin/RELEASE_BODY <<EOF
# Easel Linux Driver v$VERSION

### amd64

- Node.js v$(cat node-amd64/VERSION)
- AppImageKit v$(cat appimagekit-bin-amd64/version)

### armhf (Raspberry Pi 2/3, Zero -- ARMv7)

- Node.js v$(cat node-armhf/VERSION)
- AppImageKit v$(cat appimagekit-bin-armhf/version)

EOF
cat bin/RELEASE_BODY
