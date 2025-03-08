package utils

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/MRegterschot/GoController/models"
)

// Global includes function that checks if an element is present in the slice
func Includes(slice interface{}, value interface{}) bool {
	// Get the slice value using reflection
	v := reflect.ValueOf(slice)

	// Check if the slice is indeed a slice
	if v.Kind() != reflect.Slice {
		return false
	}

	// Iterate over the slice
	for i := 0; i < v.Len(); i++ {
		// Check if the current element in the slice is equal to the value
		if reflect.DeepEqual(v.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

// Remove function that removes an element from a slice
func Remove(slice interface{}, value interface{}) interface{} {
	// Get the slice value using reflection
	v := reflect.ValueOf(slice)

	// Check if the slice is indeed a slice
	if v.Kind() != reflect.Slice {
		return slice
	}

	// Iterate over the slice
	for i := 0; i < v.Len(); i++ {
		// Check if the current element in the slice is equal to the value
		if reflect.DeepEqual(v.Index(i).Interface(), value) {
			// Remove the element from the slice
			return reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())).Interface()
		}
	}
	return slice
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
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
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
		Items: array[start:end],
		TotalItems: len(array),
		CurrentPage: page,
		TotalPages: (len(array) + pageSize - 1) / pageSize,
		PageSize: pageSize,
	}
}

// Converts a string to an appropriate type dynamically
func ConvertStringToType(value string) interface{} {
	// Trim the input string to remove leading/trailing spaces
	value = strings.TrimSpace(value)

	// Try converting to boolean
	if value == "true" {
		return true
	} else if value == "false" {
		return false
	}

	// Try converting to integer
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	// Try converting to float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// If it's none of the above, return the string itself
	return value
}