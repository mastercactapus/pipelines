package shared

import (
	"encoding/json"
	"log"
	"os"

	"github.com/blang/semver"
)

const defaultBaseURL = "https://easel.inventables.com/downloads"

type Source struct {
	URL string
}

type Input struct {
	Source struct {
		URL string
	}
	Version struct {
		Semver *semver.Version
	}
	Params    struct{}
	OutputDir string
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
	return in
}
