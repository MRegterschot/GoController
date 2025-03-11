package utils

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/MRegterschot/GoController/models"
)

// Checks if an element is present in the slice
func Includes(slice any, value any) bool {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice {
		return false
	}

	for i := range v.Len() {
		if reflect.DeepEqual(v.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

// Removes an element from a slice.
// Returns the new slice and a boolean indicating if the element was removed
func Remove[T any](slice []T, value T) ([]T, bool) {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice {
		return slice, false
	}

	for i := range v.Len() {
		if reflect.DeepEqual(v.Index(i).Interface(), value) {
			return reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())).Interface().([]T), true
		}
	}
	return slice, false
}

// Returns the current time as an integer (seconds since the Unix epoch)
func GetCurrentTimeInSeconds() int {
	return int(time.Now().Unix())
}

// Rreturns the current time as an integer (milliseconds since the Unix epoch)
func GetCurrentTimeInMilliseconds() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

// Splits an array into chunks of a specified size
func ChunkArray[T any](array []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(array); i += chunkSize {
		end := min(i + chunkSize, len(array))
		chunks = append(chunks, array[i:end])
	}
	return chunks
}

// Paginate an array
func Paginate[T any](array []T, page int, pageSize int) models.PaginationResult[T] {
	start := page * pageSize
	end := start + pageSize
	if start > len(array) {
		start = len(array)
	}
	if end > len(array) {
		end = len(array)
	}

	return models.PaginationResult[T]{
		Items:       array[start:end],
		TotalItems:  len(array),
		CurrentPage: page,
		TotalPages:  (len(array) + pageSize - 1) / pageSize,
		PageSize:    pageSize,
	}
}

// Converts a string to an appropriate type dynamically
func ConvertStringToType(value string) any {
	value = strings.TrimSpace(value)

	if value == "true" {
		return true
	} else if value == "false" {
		return false
	}

	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	return value
}
