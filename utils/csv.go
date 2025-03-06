package utils

import (
	"encoding/csv"
	"os"

	"go.uber.org/zap"
)

// ExportCSV exports data to a CSV file
func ExportCSV(filePath string, data [][]string) error {
	err := CreateFile(filePath, nil)
	if err != nil {
		zap.L().Error("Failed to create file", zap.Error(err))
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		zap.L().Error("Failed to open file", zap.Error(err))
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(data); err != nil {
		zap.L().Error("Failed to write all", zap.Error(err))
		return err
	}

	writer.Flush()
	return writer.Error()
}