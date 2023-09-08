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
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/edgarsilva948/aftctl/pkg/aws/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// EnsureCodePipelineExists creates a new codepipeline pipeline with the given name, or returns success if it already exists.
func EnsureCodePipelineExists(client CodePipelineClient, aftManagementAccountID string, codePipelineRoleName string, pipelineName string, codeSuiteBucketName string, repoName string, branchName string, codeBuildProjectName string) (bool, error) {

	_, err := checkIfCodePipelineClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfPipelineNameIsProvided(pipelineName)

	if err != nil {
		return false, err
	}

	pipelineExists, _ := pipelineExists(client, pipelineName)

	if !pipelineExists {
		fmt.Printf("CodePipeline pipeline %s doesn't exists... creating\n", pipelineName)

		_, err := createCodePipelinePipeline(client, aftManagementAccountID, codePipelineRoleName, pipelineName, codeSuiteBucketName, repoName, branchName, codeBuildProjectName)

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

	message := fmt.Sprintf("CodePipeline Pipeline %s already exists", pipelineName)

	customCodePipelineInfoLog(logger, message)

	return true, nil
}

func customCodePipelineInfoLog(logger *zap.Logger, msg string) {

	codeEmoji := "👷 "
	coloredMsg := "\x1b[36m" + codeEmoji + " " + msg + "\x1b[0m"

	logger.Info(coloredMsg)
}

// func to create the AFT CodePipeline pipe if it doesn't exist'
func createCodePipelinePipeline(client CodePipelineClient, aftManagementAccountID string, codePipelineRoleName string, pipelineName string, codeSuiteBucketName string, repoName string, branchName string, codeBuildProjectName string) (bool, error) {

	codePipelineRoleArn := "arn:aws:iam::" + aftManagementAccountID + ":role/" + codePipelineRoleName

	input := &codepipeline.CreatePipelineInput{
		Tags: []*codepipeline.Tag{
			{
				Key:   aws.String(tags.Aftctl),
				Value: aws.String(tags.True),
			},
		},
		Pipeline: &codepipeline.PipelineDeclaration{
			Name:    aws.String(pipelineName),
			RoleArn: aws.String(codePipelineRoleArn),
			ArtifactStore: &codepipeline.ArtifactStore{
				Type:     aws.String("S3"),
				Location: aws.String(codeSuiteBucketName),
			},
			Stages: []*codepipeline.StageDeclaration{
				{
					Name: aws.String("Source"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("App"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Category: aws.String("Source"),
								Owner:    aws.String("AWS"),
								Version:  aws.String("1"),
								Provider: aws.String("CodeCommit"),
							},
							Configuration: map[string]*string{
								"RepositoryName": aws.String(repoName),
								"BranchName":     aws.String(branchName),
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{Name: aws.String("App")},
							},
							RunOrder: aws.Int64(1),
						},
					},
				},
				{
					Name: aws.String("Build"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("Build"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Category: aws.String("Build"),
								Owner:    aws.String("AWS"),
								Version:  aws.String("1"),
								Provider: aws.String("CodeBuild"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(codeBuildProjectName),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								{Name: aws.String("App")},
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{Name: aws.String("BuildOutput")},
							},
							RunOrder: aws.Int64(1),
						},
					},
				},
			},
		},
	}

	_, err := client.CreatePipeline(input)
	if err != nil {
		log.Fatalf("Error creating project: %v", err)
	}

	return true, nil
}

// func to verify if the given client is valid
func checkIfCodePipelineClientIsProvided(client CodePipelineClient) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("CodepipelineClient is not provided")
	}

	return true, nil
}

// func to verify if the given pipeline is provided
func checkIfPipelineNameIsProvided(pipelineName string) (bool, error) {
	if pipelineName == "" {
		fmt.Printf("Error: %v\n", "pipeline name is not provided")
		return false, fmt.Errorf("pipeline name is not provided")
	}

	isPipelineNameValid, err := checkPipelineNameCompliance(pipelineName)
	if !isPipelineNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given pipeline is compliant
func checkPipelineNameCompliance(pipelineName string) (bool, error) {
	length := len(pipelineName)

	// pipeline names must be between 3 (min) and 100 (max) characters long.
	if length < 3 || length > 100 {
		return false, errors.New("pipeline name must be between 3 and 100 characters long")
	}

	// pipeline names can consist only of lowercase letters, numbers, and hyphens (-).
	pattern := `^[A-Za-z0-9.@\-_]+$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(pipelineName) {
		return false, errors.New("pipeline name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number")
	}

	return true, nil
}

// pipelineExists checks if a given codebuild projct exists.
func pipelineExists(client CodePipelineClient, pipelineName string) (bool, error) {

	isPipelineExistent, err := checkIfPipelineExists(client, pipelineName)
	if err != nil {
		return false, err
	}

	return isPipelineExistent, nil
}

func checkIfPipelineExists(client CodePipelineClient, pipelineName string) (bool, error) {

	input := &codepipeline.ListPipelinesInput{}

	result, err := client.ListPipelines(input)
	if err != nil {
		return false, err
	}

	for _, existingPipeline := range result.Pipelines {
		if *existingPipeline.Name == pipelineName {
			return true, nil
		}
	}

	return false, nil
}
