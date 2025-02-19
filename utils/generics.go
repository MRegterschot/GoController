package utils

import "reflect"

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