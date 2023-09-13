/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package profile implements a new profile for AWS CLI
package profile

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	awsClient "github.com/edgarsilva948/aftctl/pkg/aws"
	"gopkg.in/ini.v1"
)

// SetupProfile creates a new profile for AWS CLI (~/.aws/credentials)
func SetupProfile(client awsClient.STSClient, accountID string, roleName string, roleSession string) error {

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)),
		RoleSessionName: aws.String(roleSession),
	}

	result, err := client.AssumeRole(input)
	if err != nil {
		return err
	}

	// Construct profile name
	profileName := fmt.Sprintf("%s-%s", accountID, roleName)

	// Get the AWS config file path
	awsConfigFile := fmt.Sprintf("%s/.aws/credentials", os.Getenv("HOME"))

	// Check if the file exists, if not, create it
	if _, err := os.Stat(awsConfigFile); os.IsNotExist(err) {
		// Create the directory if not exists
		if err := os.MkdirAll(fmt.Sprintf("%s/.aws", os.Getenv("HOME")), 0755); err != nil {
			return err
		}

		// Create the file
		_, err := os.Create(awsConfigFile)
		if err != nil {
			return err
		}
	}

	// Load the AWS config file
	cfg, err := ini.Load(awsConfigFile)
	if err != nil {
		return err
	}

	// Create a new section (AWS profile) or update an existing one
	sec, err := cfg.NewSection(profileName)
	if err != nil {
		return err
	}

	sec.Key("aws_access_key_id").SetValue(*result.Credentials.AccessKeyId)
	sec.Key("aws_secret_access_key").SetValue(*result.Credentials.SecretAccessKey)
	sec.Key("aws_session_token").SetValue(*result.Credentials.SessionToken)

	// Save the updated AWS config file
	return cfg.SaveTo(awsConfigFile)
}
