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
			return err // Return error instead of panicking
		}
		defer f.Close()

		// Format data into JSON and write it to the file
		if data != nil {
			err = json.NewEncoder(f).Encode(data)
			if err != nil {
				return err // Return error instead of panicking
			}
		}
	}
	return nil
}
