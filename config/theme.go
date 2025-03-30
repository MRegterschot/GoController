package config

import (
	"encoding/json"
	"os"

	"github.com/MRegterschot/GoController/models"
)

var Theme models.Theme

func LoadTheme() {
	theme, err := readFile[models.Theme]("./settings/theme.json")
	if err != nil {
		return
	}

	// Set the theme colors
	Theme = theme
}

// Read file and decode JSON data
func readFile[T any](file string) (T, error) {
	var data T
	f, err := os.Open(file)
	if err != nil {
		return data, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}