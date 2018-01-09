package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blang/semver"
	"github.com/mastercactapus/pipelines/easel/resource/shared"
)

func main() {
	in := shared.ReadInput()
	sel := in.Selected()

	os.MkdirAll(in.OutputDir, 0755)
	err := ioutil.WriteFile(filepath.Join(in.OutputDir, "VERSION"), []byte(sel.Semver.String()), 0644)
	if err != nil {
		log.Fatal("write VERSION file:", err)
	}

	log.Println("Download", sel.URL)
	resp, err := http.Get(sel.URL)
	if err != nil {
		log.Fatal("fetch:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("non-200:", resp.Status)
	}

	fileName := fmt.Sprintf("EaselDriver-%s.pkg", sel.Semver)
	fd, err := os.Create(filepath.Join(in.OutputDir, fileName))
	if err != nil {
		log.Fatal("create file:", err)
	}
	defer fd.Close()

	n, err := io.Copy(fd, resp.Body)
	if err != nil {
		log.Fatal("download file:", err)
	}

	type metadata struct{ Name, Value string }
	var res struct {
		Version struct {
			Semver semver.Version
		}
		Metadata []metadata
	}
	res.Version.Semver = sel.Semver
	res.Metadata = []metadata{
		{Name: "Filename", Value: fileName},
		{Name: "Size", Value: strconv.FormatInt(n, 10)},
		{Name: "Version", Value: sel.Semver.String()},
	}

	shared.WriteOutput(res)
}
