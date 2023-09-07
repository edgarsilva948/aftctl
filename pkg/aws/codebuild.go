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
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/edgarsilva948/aftctl/pkg/aws/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// EnsureCodeBuildProjectExists creates a new codebuild project with the given name, or returns success if it already exists.
func EnsureCodeBuildProjectExists(client CodeBuildClient, aftManagementAccountID string, codeBuildDockerImage string, projectName string, repoName string, repoBranch string, codeBuildRoleName string) (bool, error) {

	_, err := checkIfCodeBuildClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfProjectNameIsProvided(projectName)

	if err != nil {
		return false, err
	}

	projectExists, _ := projectExists(client, projectName)

	if !projectExists {
		fmt.Printf("CodeBuild project %s doesn't exists... creating\n", projectName)

		_, err := createCodeBuildProject(client, aftManagementAccountID, codeBuildDockerImage, projectName, repoName, repoBranch, codeBuildRoleName)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true

	logger, _ := config.Build()

	defer logger.Sync()

	message := fmt.Sprintf("CodeBuild Project %s already exists", projectName)

	customCodeBuildInfoLog(logger, message)

	return true, nil
}

func customCodeBuildInfoLog(logger *zap.Logger, msg string) {

	// CodeBuild related emojis, you can choose others
	codeEmoji := "🛠️ "
	coloredMsg := "\x1b[36m" + codeEmoji + " " + msg + "\x1b[0m"

	logger.Info(coloredMsg)
}

func checkIfProjectExists(client CodeBuildClient, projectName string) (bool, error) {
	input := &codebuild.ListProjectsInput{}

	result, err := client.ListProjects(input)
	if err != nil {
		return false, err
	}

	for _, existingProject := range result.Projects {
		if *existingProject == projectName {
			return true, nil
		}
	}

	return false, nil
}

// func to create the AFT codebuild project if it doesn't exist'
func createCodeBuildProject(client CodeBuildClient, aftManagementAccountID string, codeBuildDockerImage string, projectName string, repoName string, repoBranch string, codeBuildRoleName string) (bool, error) {

	codeBuildRoleArn := "arn:aws:iam::" + aftManagementAccountID + ":role/" + codeBuildRoleName

	input := &codebuild.CreateProjectInput{
		Tags: []*codebuild.Tag{
			{
				Key:   aws.String(tags.Aftctl),
				Value: aws.String(tags.True),
			},
		},
		Name: aws.String(projectName),
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String("CODEPIPELINE"),
		},
		Source: &codebuild.ProjectSource{
			Type: aws.String("CODEPIPELINE"),
		},
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:    aws.String("BUILD_GENERAL1_SMALL"),
			Type:           aws.String("LINUX_CONTAINER"),
			Image:          aws.String(codeBuildDockerImage),
			PrivilegedMode: aws.Bool(true),
			EnvironmentVariables: []*codebuild.EnvironmentVariable{
				{
					Name:  aws.String("REPOSITORY_NAME"),
					Value: aws.String(repoName),
				},
				{
					Name:  aws.String("REPOSITORY_BRANCH"),
					Value: aws.String(repoBranch),
				},
				// the aft deployment needs to be connected in the management account
				// https://docs.aws.amazon.com/controltower/latest/userguide/aft-getting-started.html
				{
					Name:  aws.String("MANAGEMENT_ACCOUNT_TEMPORARY_KEY_ID"),
					Value: aws.String(""),
				},
				{
					Name:  aws.String("MANAGEMENT_ACCOUNT_TEMPORARY_ACCESS_KEY"),
					Value: aws.String(""),
				},
				{
					Name:  aws.String("MANAGEMENT_ACCOUNT_TEMPORARY_SESSION_TOKEN"),
					Value: aws.String(""),
				},
			},
		},
		ServiceRole: aws.String(codeBuildRoleArn),
	}

	_, err := client.CreateProject(input)

	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}

	return true, nil
}

// func to verify if the given client is valid
func checkIfCodeBuildClientIsProvided(client CodeBuildClient) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("CodeBuildClient is not provided")
	}

	return true, nil
}

// func to verify if the given project is provided
func checkIfProjectNameIsProvided(projectName string) (bool, error) {
	if projectName == "" {
		fmt.Printf("Error: %v\n", "project name is not provided")
		return false, fmt.Errorf("project name is not provided")
	}

	isProjectNameValid, err := checkProjectNameCompliance(projectName)
	if !isProjectNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given project is compliant
func checkProjectNameCompliance(projectName string) (bool, error) {
	length := len(projectName)

	// project names must be between 3 (min) and 255 (max) characters long.
	if length < 3 || length > 255 {
		return false, errors.New("project name must be between 3 and 255 characters long")
	}

	// project names can consist only of lowercase letters, numbers, and hyphens (-).
	pattern := `^[A-Za-z0-9][A-Za-z0-9\-_]{1,254}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(projectName) {
		return false, errors.New("project name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number")
	}

	return true, nil
}

// projectExists checks if a given codebuild projct exists.
func projectExists(client CodeBuildClient, projectName string) (bool, error) {

	isProjectExistent, err := checkIfProjectExists(client, projectName)
	if err != nil {
		return false, err
	}

	return isProjectExistent, nil
}
