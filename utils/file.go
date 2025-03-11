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
