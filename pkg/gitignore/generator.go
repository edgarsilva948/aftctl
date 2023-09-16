/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package gitignore generates a .gitignore file
package gitignore

import (
	"fmt"
	"os"
	"strings"
)

// GenerateGitIgnore generates a .gitignore file
func GenerateGitIgnore() bool {
	ignoreList := []string{
		"aft-input.auto.tfvars",
		"aft-providers.tf",
		"backend.tf",
		".gitignore",
		".terraform*",
	}

	// Join the list into a single string with newlines
	content := strings.Join(ignoreList, "\n")

	// Write the content to .gitignore
	err := os.WriteFile(".gitignore", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return true

}
