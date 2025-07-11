package utils

import (
	"encoding/json"
	"os"
	"strings"
)

var translations map[string]map[string]string
var defaultLang = "es"

func LoadTranslations() error {
	translations = make(map[string]map[string]string)

	files := []string{
		"lang/es.json",
		"lang/en.json",
	}

	for _, file := range files {
		lang := strings.TrimSuffix(strings.TrimPrefix(file, "lang/"), ".json")

		data, err := os.ReadFile(file)

		if err != nil {
			return err
		}

		var m map[string]string

		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}

		translations[lang] = m
	}

	return nil
}

func Translate(key, lang string) string {
	if lang == "" {
		lang = defaultLang
	}

	if val, ok := translations[lang][key]; ok {
		return val
	}

	if val, ok := translations[defaultLang][key]; ok {
		return val
	}

	return key
}

func SetDefaultLang(lang string) {
	defaultLang = lang
}
