package utils

import (
	"reflect"
	"time"
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

// GetCurrentTimeInSeconds returns the current time as an integer (seconds since the Unix epoch)
func GetCurrentTimeInSeconds() int {
	return int(time.Now().Unix())
}

// GetCurrentTimeInMilliseconds returns the current time as an integer (milliseconds since the Unix epoch)
func GetCurrentTimeInMilliseconds() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

// ChunkArray splits an array into chunks of a specified size
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