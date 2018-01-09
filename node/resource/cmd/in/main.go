package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blang/semver"
	"github.com/mastercactapus/pipelines/node/resource/shared"
)

func main() {
	in := shared.ReadInput()

	sel := in.Selected()
	url := in.FileURL()
	file := in.FileName()

	os.MkdirAll(in.OutputDir, 0755)

	err := ioutil.WriteFile(filepath.Join(in.OutputDir, "VERSION"), []byte(in.Version.Semver.String()), 0644)
	if err != nil {
		log.Fatal("write VERSION file:", err)
	}
	fd, err := os.Create(filepath.Join(in.OutputDir, file))
	if err != nil {
		log.Fatal("create output file:", err)
	}
	defer fd.Close()

	log.Println("GET", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("download file:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("non-200 response:", resp.Status)
	}
	h1 := sha1.New()
	h256 := sha256.New()
	n, err := io.Copy(io.MultiWriter(h1, h256, fd), resp.Body)
	if err != nil {
		log.Fatal("download file:", err)
	}

	sum1 := h1.Sum(nil)
	sum256 := h256.Sum(nil)

	in.VerifyFile(sum1, sum256)

	type metadata struct{ Name, Value string }

	var res struct {
		Version struct {
			Semver semver.Version
		}
		Metadata []metadata
	}
	res.Version.Semver = *in.Version.Semver
	res.Metadata = []metadata{
		{Name: "filename", Value: file},
		{Name: "sha256", Value: hex.EncodeToString(sum256)},
		{Name: "size", Value: strconv.FormatInt(n, 10)},
		{Name: "url", Value: url},
		{Name: "semver", Value: in.Version.Semver.String()},
		{Name: "lts", Value: sel.LTS.String()},
		{Name: "versions.npm", Value: sel.NPM},
		{Name: "versions.modules", Value: sel.Modules},
		{Name: "versions.openssl", Value: sel.OpenSSL},
		{Name: "versions.uv", Value: sel.UV},
		{Name: "versions.v8", Value: sel.V8},
		{Name: "versions.zlib", Value: sel.Zlib},
	}

	shared.WriteOutput(res)
}
