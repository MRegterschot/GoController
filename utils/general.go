package utils

import (
	"encoding/base64"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

// Converts slugified base64 string back to a uuid
func DecodeSlug(slug string) (uuid.UUID, error) {
	base64Str := strings.ReplaceAll(strings.ReplaceAll(slug, "-", "+"), "_", "/") + "=="

	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.FromBytes(bytes)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}