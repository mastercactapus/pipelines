#!/bin/sh
set -e
DIR=$1
BASE=https://www.makemkv.com/download
VERSION=$(jq -r '.version.semver')

mkdir -p $DIR
cd $DIR

echo $VERSION >VERSION

curl -LsO $BASE/makemkv-sha-$VERSION.txt
gpg --verify --no-tty <makemkv-sha-$VERSION.txt

for file in $(grep '.tar.gz$' makemkv-sha-$VERSION.txt | awk '{print $2}')
do
    curl -LsO $BASE/$file
done

cat makemkv-sha-$VERSION.txt | grep '.tar.gz$' | sha256sum -c >&2
cat <<EOF
{
    "version": {"semver": "$VERSION"}
}
EOF
