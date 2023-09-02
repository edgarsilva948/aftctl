/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package helper

import (
	"reflect"
)

func HandleEscapedEmptyString(input string) string {
	if input == "\"\"" {
		input = ""
	}
	return input
}

func Contains[T comparable](slice []T, element T) bool {
	for _, sliceElement := range slice {
		if reflect.DeepEqual(sliceElement, element) {
			return true
		}
	}

	return false
}
