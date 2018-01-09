package shared

import (
	"bytes"
	"encoding/hex"
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
		Headers    bool
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
func (in Input) HeadersName() string {
	return fmt.Sprintf("node-v%s%s.tar.gz", in.Version.Semver, "-headers")
}
func (in Input) HeadersURL() string {
	return joinURL(
		in.Source.URL,
		"v"+in.Version.Semver.String()+"/"+in.HeadersName(),
	)
}
func (in Input) FileURL() string {
	return joinURL(
		in.Source.URL,
		"v"+in.Version.Semver.String()+"/"+in.FileName(),
	)
}
func (in Input) VerifyFile(file string, sum1, sum256 []byte) {
	/*
		Multiple attempts to verify the data, go as follows:
		1. SHASUMS256.txt.asc
		2. SHASUMS.txt.asc (will also use SHASUMS256.txt, if available)
		3. SHASUMS256.txt
		4. SHASUMS.txt
	*/
	once := make(map[string][]byte)
	getFile := func(name string) []byte {
		if data, ok := once[name]; ok {
			return data
		}
		url := joinURL(in.Source.URL, "v"+in.Version.Semver.String()+"/"+name)
		log.Println("GET", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("fetch %s: %+v", name, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("fetch %s: non-200 response: %s\n", name, resp.Status)
			return nil
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("read %s: %+v", name, err)
		}

		if strings.HasSuffix(name, ".asc") {
			buf := bytes.NewBuffer(data)
			cmd := exec.Command("gpg", "--verify", "--no-tty")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stderr
			cmd.Stdin = buf
			err = cmd.Run()
			if err != nil {
				log.Fatalf("failed to verify %s", name)
			}
		}
		once[name] = data
		return data
	}

	check := func(data []byte, sum []byte) bool {
		if data == nil {
			return false
		}
		lines := strings.Split(string(data), "\n")
		var hash []byte
		var err error
		for _, l := range lines {
			if strings.HasSuffix(l, file) {
				hash, err = hex.DecodeString(strings.SplitN(l, " ", 2)[0])
				if err != nil {
					log.Fatal("decode checksum:", err)
				}
				break
			}
		}
		if hash == nil {
			return false
		}
		if !bytes.Equal(hash, sum) {
			log.Fatalf("checksum mismatch: got '%s' but expected '%s'", hex.EncodeToString(sum), hex.EncodeToString(hash))
		}
		return true
	}

	// if check(getFile("SHASUMS256.txt.asc"), sum256) {
	// 	return
	// }
	// if check(getFile("SHASUMS.txt.asc"), sum1) {
	// 	// don't care if the file exists, but we validate if it is there
	// 	check(getFile("SHASUMS256.txt"), sum256)
	// 	return
	// }
	if check(getFile("SHASUMS256.txt"), sum256) {
		return
	}
	if check(getFile("SHASUMS.txt"), sum1) {
		return
	}
	log.Println("WARNING: Could not find any checksum to validate file!!!")
	return
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
