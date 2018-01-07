#!/bin/sh
set -ex

OUT=$(pwd)/bin
mkdir -p $OUT

GO_REPO=github.com/mastercactapus/pipelines
GO_PKG=$GO_REPO/easel/resource

mkdir -p $(go env GOPATH)/src/$GO_REPO
cp -r $(pwd)/repo/. $(go env GOPATH)/src/$GO_REPO/
cp repo/easel/resource/Dockerfile $OUT/

export CGO_ENABLED=0
go get $GO_PKG/...
go build -o $OUT/check $GO_PKG/cmd/check
go build -o $OUT/in $GO_PKG/cmd/in
