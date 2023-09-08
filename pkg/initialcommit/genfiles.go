/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package initialcommit creates the repo and give the user instructions to push
package initialcommit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GenerateCommitFiles creates the directory and files to be pushed
func GenerateCommitFiles(repoName string, tfBucket string, region string, tfVersion string) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true

	fileEmoji := "📄"
	dirEmoji := "📁"

	logger, _ := config.Build()

	defer logger.Sync()

	// creating the dir with the repo name
	_, err := ensureDirExists(repoName)
	if err != nil {
		log.Fatalf("Error creating the directory: %v", err)
	}

	dirMessage := "directory file successfully generated"
	customCreateFilesLog(logger, dirMessage, dirEmoji)

	// creating the backend.tf file
	_, err = createBackendtfFile(filepath.Join(repoName, "backend.tf"), tfBucket, region)
	if err != nil {
		log.Fatalf("Error creating the backend.tf file: %v", err)
	}

	backendtfMessage := "backend.tf file successfully generated"
	customCreateFilesLog(logger, backendtfMessage, fileEmoji)

	// creating the buildspec.yaml file
	_, err = createBuildSpecFile(filepath.Join(repoName, "buildspec.yaml"), tfVersion)
	if err != nil {
		log.Fatalf("Error creating the buildspec.yaml file: %v", err)
	}

	buildSpecMessage := "buildspec.yaml file successfully generated"
	customCreateFilesLog(logger, buildSpecMessage, fileEmoji)

}

func customCreateFilesLog(logger *zap.Logger, msg string, emoji string) {

	coloredMsg := "\x1b[32m" + emoji + " " + msg + "\x1b[0m"

	logger.Info(coloredMsg)
}

func createBackendtfFile(path string, tfBucket string, region string) (bool, error) {
	// Construct the content for the backend.tf file
	content := fmt.Sprintf(`terraform {
	backend "s3" {
		bucket = "%s"
		key    = "tfstate"
		region = "%s"
	}
}`, tfBucket, region)

	// Write the content to the backend.tf file
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to backend.tf: %s", err)
		return false, err
	}

	return true, nil
}

func createBuildSpecFile(path string, tfVersion string) (bool, error) {

	content := fmt.Sprintf(`version: 0.2
env:
  variables:
    TERRAFORM_VER: "%s"
phases:
  install:
    commands:
      - export TERRAFORM_VER=$TERRAFORM_VER
      - export AWS_ACCESS_KEY_ID=$MANAGEMENT_ACCOUNT_TEMPORARY_KEY_ID
      - export AWS_SECRET_ACCESS_KEY=$MANAGEMENT_ACCOUNT_TEMPORARY_ACCESS_KEY
      - export AWS_SESSION_TOKEN=$MANAGEMENT_ACCOUNT_TEMPORARY_SESSION_TOKEN
      - echo $TERRAFORM_VER
      - yum install -y wget unzip
      - wget -q https://releases.hashicorp.com/terraform/${TERRAFORM_VER}/terraform_${TERRAFORM_VER}_linux_amd64.zip
      - unzip terraform_${TERRAFORM_VER}_linux_amd64.zip
      - mv terraform /usr/local/bin/
  build:
    commands:
      - terraform init
      - terraform plan
      - terraform apply --auto-approve
  post_build:
    commands:
      - echo "AFT deployment successfully" > definitions.json
artifacts:
  files: definitions.json`, tfVersion)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to buildspec.yaml: %s", err)
		return false, err
	}

	return true, nil
}

func ensureDirExists(dir string) (bool, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err == nil {
			fmt.Printf("Directory %s created successfully.\n", dir)
			return true, nil
		}

		log.Fatalf("Failed to create directory: %s", err)
		return false, err
	} else if err == nil {
		fmt.Printf("Directory %s already exists. Skipping creation.\n", dir)
		return true, nil
	} else {
		log.Fatalf("Error checking directory: %s", err)
		return false, err
	}
}
