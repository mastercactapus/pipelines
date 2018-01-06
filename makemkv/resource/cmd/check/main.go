package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

const baseURL = "https://www.makemkv.com/download/"

func main() {
	c := colly.NewCollector()
	var latest string
	rx := regexp.MustCompile(`/download/makemkv-sha-([0-9]+\.[0-9]+\.[0-9]+)\.txt$`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		m := rx.FindStringSubmatch(e.Attr("href"))
		if m == nil {
			return
		}
		latest = m[1]
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println(string(r.Body))
		panic(err)
	})

	c.Visit(baseURL)

	type version struct {
		Semver string `json:"semver"`
	}

	out := []version{{Semver: latest}}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(out)
	if err != nil {
		panic(err)
	}
}
