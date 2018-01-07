package shared

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
)

type LTS string

func (l *LTS) UnmarshalJSON(b []byte) error {
	if string(b) == "false" {
		*l = ""
		return nil
	}
	if string(b) == "true" {
		*l = "*"
		return nil
	}
	var s string
	err := json.Unmarshal(b, &s)
	*l = LTS(s)
	return err
}
func (l LTS) Match(b LTS) bool {
	if l == "" {
		return true
	}
	if l == b {
		return true
	}
	if l == "*" && b != "" {
		return true
	}
	return false
}

type Version struct {
	Semver     semver.Version
	RawVersion string `json:"version"`
	Files      []string
	Date       string

	NPM     string
	V8      string
	UV      string
	Zlib    string
	OpenSSL string
	Modules string
	LTS     LTS
}

var rx = regexp.MustCompile(`\b([0-9]+\.[0-9]+\.[0-9]+)\b.*\.pkg\b`)

func (in Input) Selected() Version {
	v := in.LatestVersions()
	if len(v) == 0 {
		log.Fatal("Version not found: %s", in.Version.Semver)
	}

	return v[0]
}

func joinURL(a, b string) string {
	return strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(b, "/")
}

func contains(strs []string, s string) bool {
	for _, chk := range strs {
		if chk == s {
			return true
		}
	}
	return false
}

func (in Input) LatestVersions() []Version {
	log.Println("Fetching available versions...")
	resp, err := http.Get(joinURL(in.Source.URL, "index.json"))
	if err != nil {
		log.Fatal("fetch index.json:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("non-200 response:", resp.Status)
	}

	var versions []Version
	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		log.Fatal("decode:", err)
	}

	filtered := versions[:0]
	for _, v := range versions {
		v.Semver = semver.MustParse(strings.TrimPrefix(v.RawVersion, "v"))
		if !in.Source.LTS.Match(v.LTS) {
			continue
		}
		if in.Version.Semver != nil && in.Version.Semver.GT(v.Semver) {
			continue
		}
		if !in.Range(v.Semver) {
			continue
		}
		if in.Params.Classifier != "" && !contains(v.Files, in.Params.Classifier) {
			continue
		}
		filtered = append(filtered, v)
	}
	versions = filtered

	sort.Slice(versions, func(i, j int) bool { return versions[i].Semver.LT(versions[j].Semver) })

	if in.Version.Semver == nil {
		return []Version{versions[len(versions)-1]}
	}

	return versions
}
