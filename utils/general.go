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

// GetCurrentTimeInSeconds returns the current time as an integer (seconds since the Unix epoch)
func GetCurrentTimeInSeconds() int {
	return int(time.Now().Unix())
}

// GetCurrentTimeInMilliseconds returns the current time as an integer (milliseconds since the Unix epoch)
func GetCurrentTimeInMilliseconds() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}