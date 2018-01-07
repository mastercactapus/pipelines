package shared

import (
	"encoding/json"
	"log"
	"os"
)

func WriteOutput(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(v)
	if err != nil {
		log.Fatalln("encode output:", err)
	}
}
