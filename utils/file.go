package utils

import (
	"encoding/json"
	"os"
)

// Create file with data if it doesn't exist
func CreateFile(file string, data any) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()

		if data != nil {
			err = json.NewEncoder(f).Encode(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Read file and decode JSON data
func ReadFile[T any](file string) (T, error) {
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