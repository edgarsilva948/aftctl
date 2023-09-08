/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package initialcommit creates the repo and give the user instructions to push
package initialcommit

import (
	"fmt"
	"strings"
)

// PushCode function is used to show push instructions
func PushCode(repoPath string, region string, repoName string) {

	interpolatedRepoURL := "https://git-codecommit." + region + ".amazonaws.com/v1/repos/" + repoName

	commands := [][]string{
		{"git", "init"},
		{"git", "remote", "add", "origin", interpolatedRepoURL},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "Initial commit"},
		{"git", "push", "origin", "main"},
	}

	// Instructions for checking if 'main' branch exists
	fmt.Println("First, check if a 'main' branch exists:")
	fmt.Printf("  git rev-parse --verify --quiet main\n")
	fmt.Println("If the above command does not produce any output, it means the 'main' branch doesn't exist. Create it using the command below:")
	fmt.Println("  git checkout -b main")

	// Print each command in sequence
	fmt.Println("Then, execute the following commands in your terminal:")
	for _, command := range commands {
		fmt.Printf("  %s\n", joinArgs(command))
	}
}

// Helper function to join command and its arguments
func joinArgs(args []string) string {
	return strings.Join(args, " ")
}
