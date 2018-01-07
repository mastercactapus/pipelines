package shared

import (
	"log"
	"regexp"
	"sort"

	"github.com/blang/semver"
	"github.com/gocolly/colly"
)

type Version struct {
	Semver semver.Version
	URL    string
}

var rx = regexp.MustCompile(`\b([0-9]+\.[0-9]+\.[0-9]+)\b.*\.pkg\b`)

func (in Input) Selected() Version {
	v := in.LatestVersions()
	if len(v) == 0 {
		log.Fatal("Version not found: %s", in.Version.Semver)
	}

	return v[0]
}

func (in Input) LatestVersions() []Version {
	log.Println("Fetching available versions...")
	c := colly.NewCollector()
	var versions []Version
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		m := rx.FindStringSubmatch(e.Attr("href"))
		if m == nil {
			return
		}

		versions = append(versions, Version{Semver: semver.MustParse(m[1]), URL: e.Attr("href")})
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println(string(r.Body))
		panic(err)
	})

	c.Visit(in.Source.URL)

	if len(versions) == 0 {
		return nil
	}

	filtered := versions[:0]
	for _, v := range versions {
		if in.Version.Semver.GT(v.Semver) {
			continue
		}
		filtered = append(filtered, v)
	}
	versions = filtered

	sort.Slice(versions, func(i, j int) bool { return versions[i].Semver.LT(versions[j].Semver) })

	if in.Version.Semver == nil {
		return []Version{versions[len(versions)-1]}
	}

	log.Println("Listing versions since", in.Version.Semver)

	return versions
}
