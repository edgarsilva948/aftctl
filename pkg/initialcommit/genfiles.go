/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package initialcommit creates the repo and give the user instructions to push
package initialcommit

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/edgarsilva948/aftctl/pkg/logging"
)

const fileEmoji = "📄"
const dirEmoji = "📁"
const zipEmoji = "📦"

// GenerateCommitFiles creates the directory and files to be pushed
func GenerateCommitFiles(
	repoName string,
	tfBucket string,
	region string,
	tfVersion string,
	ctManagementAccountID string,
	logArchiveAccountID string,
	auditAccountID string,
	aftManagementAccountID string,
	ctHomeRegion string,
	tfBackendSecondaryRegion string,
	aftMetricsReporting bool,
	aftFeatureCloudtrailDataEvents bool,
	aftFeatureEnterpriseSupport bool,
	aftFeatureDeleteDefaultVPCsEnabled bool,
	terraformDistribution string,
) {

	// creating the dir with the repo name
	message, color, err := ensureDirExists(repoName, dirEmoji)
	if err != nil {
		log.Fatalf("Error creating the directory: %v", err)
	}

	logging.CustomLog(dirEmoji, color, message)

	// creating the backend.tf file
	message, err = createBackendtfFile(repoName, fileEmoji, tfBucket, region)

	if err != nil {
		log.Fatalf("Error creating the backend.tf file: %v", err)
	}

	logging.CustomLog(fileEmoji, "green", message)

	// creating the buildspec.yaml file
	message, err = createBuildSpecFile(repoName, fileEmoji, tfVersion)
	if err != nil {
		log.Fatalf("Error creating the buildspec.yaml file: %v", err)
	}

	logging.CustomLog(fileEmoji, "green", message)

	// creating the main.tf file
	message, err = createMainTFFile(repoName,
		fileEmoji,
		ctManagementAccountID,
		logArchiveAccountID,
		auditAccountID,
		aftManagementAccountID,
		ctHomeRegion,
		tfBackendSecondaryRegion,
		aftMetricsReporting,
		aftFeatureCloudtrailDataEvents,
		aftFeatureEnterpriseSupport,
		aftFeatureDeleteDefaultVPCsEnabled,
		tfVersion,
		terraformDistribution,
	)

	if err != nil {
		log.Fatalf("Error creating the main.tf file: %v", err)
	}

	logging.CustomLog(fileEmoji, "green", message)

	message, err = zipDirectory(repoName, zipEmoji)
	if err != nil {
		fmt.Println("Error creating the zip file:", err)
	}

	logging.CustomLog(zipEmoji, "green", message)

}

func createBackendtfFile(dir string, fileEmoji string, tfBucket string, region string) (string, error) {

	content := fmt.Sprintf(`terraform {
	backend "s3" {
		bucket = "%s"
		key    = "tfstate"
		region = "%s"
	}
}`, tfBucket, region)

	path := filepath.Join(dir, "backend.tf")

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "Failed to write to buildspec.yaml", err
	}

	message := "File ./" + dir + "/backend.tf successfully created"
	return message, nil
}

func createBuildSpecFile(dir string, fileEmoji string, tfVersion string) (string, error) {

	content := fmt.Sprintf(`version: 0.2
env:
  variables:
    TERRAFORM_VERSION: "%s"
phases:
  install:
    commands:
	  - |
        set -e
		echo $TERRAFORM_VERSION
		echo "Installing terraform"
		cd /tmp
		curl -q -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
		unzip -q -o terraform_${TERRAFORM_VERSION}_linux_amd64.zip
		mv terraform /usr/local/bin/
		terraform -no-color --version
  build:
    on-failure: ABORT
	commands:
	  - |
	    cd $CODEBUILD_SRC_DIR/terraform
        echo "Running terraform apply"
        terraform apply -no-color -input=false --auto-approve "output.tfplan"
  post_build:
    commands:
      - echo "AFT setup deployment successfully"
artifacts:
  files:
    - '**/*'
`, tfVersion)

	path := filepath.Join(dir, "buildspec.yaml")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "Failed to write to buildspec.yaml", err
	}

	message := "File ./" + dir + "/buildspec.yaml successfully created"
	return message, nil
}

func ensureDirExists(dir string, dirEmoji string) (string, string, error) {
	var err error
	if _, err := os.Stat(dir); err == nil {
		color := "blue"
		message := "Directory " + dir + " already exists"
		return message, color, nil

	} else if os.IsNotExist(err) {
		color := "green"
		message := "Directory " + dir + " successfully created"

		return message, color, os.Mkdir(dir, 0755)
	}

	return "Error creating the repo directory", "red", err
}

func createMainTFFile(
	dir string,
	fileEmoji string,
	ctManagementAccountID string,
	logArchiveAccountID string,
	auditAccountID string,
	aftManagementAccountID string,
	ctHomeRegion string,
	tfBackendSecondaryRegion string,
	aftMetricsReporting bool,
	aftFeatureCloudtrailDataEvents bool,
	aftFeatureEnterpriseSupport bool,
	aftFeatureDeleteDefaultVPCsEnabled bool,
	terraformVersion string,
	terraformDistribution string,
) (string, error) {

	aftDeployTemplate := fmt.Sprintf(`
# Copyright Amazon.com, Inc. or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

module "aft" {

  source = "github.com/aws-ia/terraform-aws-control_tower_account_factory"
  
  # Required variables
  ct_management_account_id  = "%s"
  log_archive_account_id    = "%s"
  audit_account_id          = "%s"
  aft_management_account_id = "%s"
  ct_home_region            = "%s"
 
  # Optional variables
  tf_backend_secondary_region = "%s"
  aft_metrics_reporting       = "%t"
 
  # AFT Feature flags
  aft_feature_cloudtrail_data_events      = "%t"
  aft_feature_enterprise_support          = "%t"
  aft_feature_delete_default_vpcs_enabled = "%t"
 
  # Terraform variables
  terraform_version      = "%s"
  terraform_distribution = "%s"
}`,
		ctManagementAccountID, logArchiveAccountID, auditAccountID, aftManagementAccountID, ctHomeRegion,
		tfBackendSecondaryRegion, aftMetricsReporting, aftFeatureCloudtrailDataEvents, aftFeatureEnterpriseSupport,
		aftFeatureDeleteDefaultVPCsEnabled, terraformVersion, terraformDistribution)

	// check if management account id is valid
	validAccount, err := isValidAWSAccountID(ctManagementAccountID)
	if !validAccount {
		err := fmt.Errorf("management account ID is not valid: %v", err)
		log.Println(err)
	}

	// check if log archive account id is valid
	validAccount, err = isValidAWSAccountID(logArchiveAccountID)
	if !validAccount {
		err := fmt.Errorf("log Archive account ID is not valid: %v", err)
		log.Println(err)
	}

	// check if audit account id is valid
	validAccount, err = isValidAWSAccountID(auditAccountID)
	if !validAccount {
		err := fmt.Errorf("audit account ID is not valid: %v", err)
		log.Println(err)
	}

	// check if aft account id is valid
	validAccount, err = isValidAWSAccountID(aftManagementAccountID)
	if !validAccount {
		err := fmt.Errorf("AFT account ID is not valid: %v", err)
		log.Println(err)
	}

	path := filepath.Join(dir, "main.tf")
	err = os.WriteFile(path, []byte(aftDeployTemplate), 0644)
	if err != nil {
		return "Failed to write to main.tf", err
	}

	message := "File ./" + dir + "/main.tf successfully created"
	return message, nil

}

// isValidAWSAccountID checks if a string represents a valid account id
func isValidAWSAccountID(accountID string) (bool, error) {
	var err error
	if len(accountID) != 12 {
		return false, err
	}

	// Check if all characters are digits
	_, err = strconv.ParseUint(accountID, 10, 64)
	if err != nil {
		return false, err
	}

	return true, nil
}

func zipDirectory(dir string, fileEmoji string) (string, error) {
	// Check if directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "Directory does not exist.", err
	}
	if !info.IsDir() {
		return "Provided path is not a directory.", fmt.Errorf("not a directory: %s", dir)
	}

	// Create zip file
	zipFileName := dir + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return "An error occurred when generating the zip file.", err
	}
	defer zipFile.Close()

	// Initialize a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Variable to hold any error that occurs during filepath.Walk
	var walkErr error

	// Walk through each file in the directory
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkErr = err
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			walkErr = err
			return err
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			walkErr = err
			return err
		}

		fsFile, err := os.Open(path)
		if err != nil {
			walkErr = err
			return err
		}
		defer fsFile.Close()

		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			walkErr = err
		}
		return err
	})

	// Check if any error occurred during filepath.Walk
	if walkErr != nil {
		return "An error occurred while zipping the directory.", walkErr
	}

	message := "File ./" + dir + ".zip successfully created"
	return message, nil
}
