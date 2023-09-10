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
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/edgarsilva948/aftctl/pkg/aws/tags"
	"github.com/edgarsilva948/aftctl/pkg/logging"
)

const cfnIcon = "📚"

// EnsureCloudformationExists creates a new cloudformation stack with the given name, or returns success if it already exists.
func EnsureCloudformationExists(client CloudformationClient, stackName string, repoName string, repoDescription string, bucketName string, zipFileName string) (bool, error) {

	_, err := checkIfCloudformationClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfStackNameIsProvided(stackName)

	if err != nil {
		return false, err
	}

	stackExists, _ := stackExists(client, stackName)

	if !stackExists {

		message := fmt.Sprintf("Cloudformation stack %s doesn't exists... creating", stackName)
		logging.CustomLog(cfnIcon, "yellow", message)

		_, err := createStack(client, stackName, repoName, repoDescription, bucketName, zipFileName)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	message := fmt.Sprintf("Cloudformation Stack %s already exists", stackName)
	logging.CustomLog(cfnIcon, "blue", message)

	return true, nil

}

// func to verify if the given client is valid
func checkIfCloudformationClientIsProvided(client CloudformationClient) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("CloudformationClient is not provided")
	}

	return true, nil
}

// func to verify if the given stack name is provided
func checkIfStackNameIsProvided(stackName string) (bool, error) {
	if stackName == "" {
		fmt.Printf("Error: %v\n", "stack name is not provided")
		return false, fmt.Errorf("stack name is not provided")
	}

	isStackNameValid, err := checkStackNameCompliance(stackName)
	if !isStackNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given stack name is compliant
func checkStackNameCompliance(stackName string) (bool, error) {
	length := len(stackName)

	// stack names must be between 3 (min) and 100 (max) characters long.
	if length < 1 || length > 100 {
		return false, errors.New("stack name must be between 1 and 100 characters long")
	}

	pattern := `^[a-zA-Z0-9-_]+$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(stackName) {
		return false, errors.New("stack name can only consist of lowercase letters, numbers, and hyphens")
	}
	return true, nil
}

// stackExists checks if a given cloudformation stack exists.
func stackExists(client CloudformationClient, stackName string) (bool, error) {

	isStackExistent, err := checkIfStackExists(client, stackName)
	if err != nil {
		return false, err
	}

	return isStackExistent, nil
}

// func to verify if the given stack name already exists
func checkIfStackExists(client CloudformationClient, stackName string) (bool, error) {

	input := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	_, err := client.DescribeStacks(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

// func to create given stack if it doesn't exist'
func createStack(client CloudformationClient, stackName string, repoName string, repoDescription string, bucketName string, zipFileName string) (bool, error) {

	// Define CloudFormation template
	template := `
Resources:
  MyCodeCommitRepository:
    Type: "AWS::CodeCommit::Repository"
    Properties:
      RepositoryName: "%s"
      RepositoryDescription: "%s"
      Tags:
        - Key: "%s"
          Value: "%s"
      Code:
        S3:
          Bucket: "%s"
          Key: "%s"
`

	input := &cloudformation.CreateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(fmt.Sprintf(template, repoName, repoDescription, tags.Aftctl, tags.True, bucketName, zipFileName)),
	}

	_, err := client.CreateStack(input)
	if err != nil {
		log.Fatalf("Error creating CloudFormation stack: %v", err)
	}

	message := fmt.Sprintf("Cloudformation stack %s successfully created", stackName)
	logging.CustomLog(secIcon, "green", message)

	return true, nil
}
