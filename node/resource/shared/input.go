package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/blang/semver"
)

const defaultBaseURL = "https://nodejs.org/dist"

type Source struct {
	URL string
}

type Input struct {
	Source struct {
		URL    string
		LTS    LTS
		Range  string
		_range semver.Range
	}
	Version struct {
		Semver *semver.Version
	}
	Params struct {
		Classifier string
	}
	OutputDir string
}

func (in Input) FileName() string {
	spec := ""
	if in.Params.Classifier == "src" {
		// nothing
	} else if in.Params.Classifier != "" {
		spec = "-" + in.Params.Classifier
	}
	return fmt.Sprintf("node-v%s%s.tar.gz", in.Version.Semver, spec)
}
func (in Input) FileURL() string {
	return joinURL(
		in.Source.URL,
		"v"+in.Version.Semver.String()+"/"+in.FileName(),
	)
}
func (in Input) FileSHA256() string {
	url := joinURL(in.Source.URL, "v"+in.Version.Semver.String()+"/SHASUMS256.txt.asc")
	log.Println("GET", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("fetch SHASUMS256.txt.asc:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("non-200 response:", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read SHASUMS256.txt.asc:", err)
	}

	buf := bytes.NewBuffer(data)
	cmd := exec.Command("gpg", "--verify", "--no-tty")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	cmd.Stdin = buf
	err = cmd.Run()
	if err != nil {
		log.Fatal("failed to verify SHASUMS256.txt.asc")
	}
	lines := strings.Split(string(data), "\n")
	file := in.FileName()
	for _, l := range lines {
		if strings.HasSuffix(l, file) {
			return strings.SplitN(l, " ", 2)[0]
		}
	}

	log.Fatalf("failed to find entry for '%s' in SHASUMS256.txt.asc")
	return ""
}

func (in Input) Range(v semver.Version) bool {
	if in.Source._range == nil {
		return true
	}
	return in.Source._range(v)
}

func ReadInput() Input {
	var in Input
	err := json.NewDecoder(os.Stdin).Decode(&in)
	if err != nil {
		log.Fatal("parse input JSON:", err)
	}
	if in.Source.URL == "" {
		in.Source.URL = defaultBaseURL
	}
	if len(os.Args) == 2 {
		in.OutputDir = os.Args[1]
	} else {
		in.OutputDir = "."
	}
	if in.Source.Range != "" {
		in.Source._range, err = semver.ParseRange(in.Source.Range)
		if err != nil {
			log.Fatal("parse source range:", err)
		}
	}
	return in
}
