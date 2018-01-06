#!/bin/sh
set -ex

OUT=$(pwd)/bin
mkdir -p $OUT

cp -r $(pwd)/repo/makemkv/resource/cmd/check $(go env GOPATH)/src/check
cp repo/makemkv/resource/cmd/in/in.sh $OUT/in
cp repo/makemkv/resource/Dockerfile $OUT/

go get check/...
go build -o $OUT/check check
