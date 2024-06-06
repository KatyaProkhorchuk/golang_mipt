package configs

import (
	"encoding/json"
	"log"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal"
)

func GetLanguages(name string) []internal.Language {
	var languages []internal.Language
	file, err := os.ReadFile(name)

	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &languages)
	if err != nil {
		log.Fatal(err)
	}
	return languages
}
