/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aws

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/edgarsilva948/aftctl/pkg/aws/tags"
)

// EnsureCodeCommitRepoExists creates a new codecommit repository with the given name, or returns success if it already exists.
func EnsureCodeCommitRepoExists(client CodeCommitClient, repoName string, description string) (bool, error) {

	_, err := checkIfCodeCommitClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfRepoNameIsProvided(repoName)

	if err != nil {
		return false, err
	}

	repoExists, _ := repoExists(client, repoName)

	if !repoExists {
		fmt.Printf("CodeCommit repository %s doesn't exists... creating\n", repoName)

		_, err := createRepo(client, repoName, description)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	fmt.Printf("CodeCommit repository %s already exists... continuing", repoName)

	return true, nil
}

// func to verify if the given repository is provided
func checkIfRepoNameIsProvided(repoName string) (bool, error) {
	if repoName == "" {
		fmt.Printf("Error: %v\n", "repository name is not provided")
		return false, fmt.Errorf("repository name is not provided")
	}

	isRepoNameValid, err := checkRepoNameCompliance(repoName)
	if !isRepoNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given client is valid
func checkIfCodeCommitClientIsProvided(client CodeCommitClient) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("CodeCommitClient is not provided")
	}

	return true, nil
}

// func to verify if the given repository is compliant
func checkRepoNameCompliance(repoName string) (bool, error) {
	length := len(repoName)

	// repository names must be between 3 (min) and 63 (max) characters long.
	if length < 1 || length > 100 {
		return false, errors.New("repository name must be between 1 and 100 characters long")
	}

	pattern := `^[a-zA-Z0-9-_]+$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(repoName) {
		return false, errors.New("repository name can only consist of lowercase letters, numbers, and hyphens")
	}
	return true, nil
}

// repoExists checks if a given codecommit repo exists.
func repoExists(client CodeCommitClient, repoName string) (bool, error) {

	isRepoExistent, err := checkIfRepoExists(client, repoName)
	if err != nil {
		return false, err
	}

	return isRepoExistent, nil
}

// func to verify if the given repo name already exists
func checkIfRepoExists(client CodeCommitClient, repoName string) (bool, error) {
	input := &codecommit.GetRepositoryInput{
		RepositoryName: aws.String(repoName),
	}

	_, err := client.GetRepository(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

// func to create given repo if it doesn't exist'
func createRepo(client CodeCommitClient, repoName string, description string) (bool, error) {

	_, err := client.CreateRepository(&codecommit.CreateRepositoryInput{
		RepositoryName:        aws.String(repoName),
		RepositoryDescription: aws.String(description),
		Tags: map[string]*string{
			tags.Aftctl: aws.String(tags.True),
		},
	})

	if err != nil {
		log.Printf("unable to create repository %q, %v", repoName, err)
		return false, err
	}

	fmt.Printf("CodeCommit repository %s successfuly created", repoName)
	return true, nil
}
