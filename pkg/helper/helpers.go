/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package helper

import (
	"reflect"
)

// HandleEscapedEmptyString takes an input string and returns an empty string if the input is "\"\"".
func HandleEscapedEmptyString(input string) string {
	if input == "\"\"" {
		input = ""
	}
	return input
}

// Contains checks if a slice contains a given element.
func Contains(slice []interface{}, element interface{}) bool {
	for _, sliceElement := range slice {
		if reflect.DeepEqual(sliceElement, element) {
			return true
		}
	}

	return false
}
