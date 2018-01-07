package main

import (
	"github.com/blang/semver"
	"github.com/mastercactapus/pipelines/easel/resource/shared"
)

func main() {
	type out struct {
		Semver semver.Version
	}

	var res []out
	for _, v := range shared.ReadInput().LatestVersions() {
		res = append(res, out{Semver: v.Semver})
	}

	shared.WriteOutput(res)
}
