/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package local provides the local command
package local

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"os"

	"path/filepath"
	"strings"

	"github.com/edgarsilva948/aftctl/pkg/aws"
	profile "github.com/edgarsilva948/aftctl/pkg/aws/profiles"
	"github.com/edgarsilva948/aftctl/pkg/gitignore"
	validate "github.com/edgarsilva948/aftctl/pkg/validator"
	"github.com/flosch/pongo2"

	"github.com/caarlos0/log"

	"github.com/spf13/cobra"
)

const (
	aftMgmtAccountID     = "/aft/account/aft-management/account-id"
	tfDistribution       = "/aft/config/terraform/distribution"
	ctMgmtRegion         = "/aft/config/ct-management-region"
	aftAdminRoleName     = "/aft/resources/iam/aft-administrator-role-name"
	aftExecutionRoleName = "/aft/resources/iam/aft-execution-role-name"
	aftSessionName       = "/aft/resources/iam/aft-session-name"
	tfBackendRegion      = "/aft/config/oss-backend/primary-region"
	tfKmsKeyID           = "/aft/config/oss-backend/kms-key-id"
	tfDynamoDBTableName  = "/aft/config/oss-backend/table-id"
	tfS3BucketID         = "/aft/config/oss-backend/bucket-id"
	tfVarFile            = "aft-input.auto.tfvars"
)

var args struct {
	targetAccount    string
	terraformCommand string
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false

	flags.StringVarP(
		&args.targetAccount,
		"target-account",
		"a",
		"",
		"Account ID to be targeted during local execution",
	)

	flags.StringVarP(
		&args.terraformCommand,
		"terraform-command",
		"c",
		"",
		"The terraform command to be executed locally",
	)
}

// Cmd represents the Cobra command for the local AFT execution.
var Cmd = &cobra.Command{
	Use:   "local",
	Short: "Runs AFT locally",
	Long:  "Runs AFT locally executing the same commands as the pipeline",
	Run:   Run,
}

// Run executes the local command
func Run(cmd *cobra.Command, argv []string) {

	// client initialization with AFT Credentials
	awsClient, ssmClient, err := initializeAWSandSSMClients()
	if err != nil {
		log.Errorf("error initializing AWS and SSM Clients: %v", err)
		return
	}

	// defining the S3 key for the local execution
	var tfS3Key string

	// Validate input
	if err := validateInput(); err != nil {
		log.Errorf("Validation failed: %v", err)
		return
	}

	// Define an array of SSM parameter keys that we need to fetch.
	ssmKeys := []string{
		aftMgmtAccountID,
		tfDistribution,
		ctMgmtRegion,
		aftAdminRoleName,
		aftExecutionRoleName,
		tfBackendRegion,
		tfKmsKeyID,
		tfDynamoDBTableName,
		tfS3BucketID,
	}

	// Fetch the SSM parameters based on the keys defined above.
	params := getSSMParameters(ssmClient, ssmKeys)

	// Check for DynamoDB table name parameter.
	tfDynamoDBTableNameParam := params[tfDynamoDBTableName]

	// Check for Management Account ID parameter.
	aftMgmtAccountIDParam := params[aftMgmtAccountID]

	// Check for Admin Role Name parameter.
	aftAdminRoleNameParam := params[aftAdminRoleName]

	// Check for Distribution Type parameter.
	tfDistributionParam := params[tfDistribution]

	// Check for Management Region parameter.
	ctMgmtRegionParam := params[ctMgmtRegion]

	// Check for Backend Region parameter.
	tfBackendRegionParam := params[tfBackendRegion]

	// Check for Execution Role Name parameter.
	aftExecutionRoleNameParam := params[aftExecutionRoleName]

	// Check for S3 Bucket ID parameter.
	tfS3BucketIDParam := params[tfS3BucketID]

	// Check for KMS Key ID parameter.
	tfKmsKeyIDParam := params[tfKmsKeyID]

	// getting the current directory
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory:", err)
		return
	}

	// Call the getTFS3Key function to determine the appropriate S3 key based on the current directory.
	tfS3Key, err = getTFS3Key(pwd, args.targetAccount)
	if err != nil {
		return
	}

	// setting up the AWS profile for the AFT Account using the user current credentials
	if err := profile.SetupProfile(awsClient.GetSTSClient(), aftMgmtAccountIDParam, aftAdminRoleNameParam, "AWSAFT-Session"); err != nil {
		fmt.Println("error setting up profile", err)
	}

	log.Infof("successfully set up AWS profile %s-%s", aftMgmtAccountIDParam, aftAdminRoleNameParam)

	// Assuming the AFT Admin Role in the AFT Account
	profileVariable := aftMgmtAccountIDParam + "-" + aftAdminRoleNameParam
	aws.NewClient(profileVariable)
	accessKey, secretKey, sessionToken, err := aws.GetAWSCredentials(profileVariable)
	if err != nil {
		fmt.Println("error getting AFT Admin credentials")
	}

	// Set the AWS_PROFILE environment variable
	err = os.Setenv("AWS_PROFILE", profileVariable)
	if err != nil {
		log.Fatalf("error setting env var AWS_PROFILE: %v", err)
	}

	// calling the function to process Jinja files
	processJinjaFiles(
		tfDistributionParam,
		ctMgmtRegionParam,
		tfBackendRegionParam,
		aftMgmtAccountIDParam,
		aftExecutionRoleNameParam,
		tfS3BucketIDParam,
		tfS3Key,
		tfDynamoDBTableNameParam,
		tfKmsKeyIDParam,
	)

	// Generate the .gitignore file
	gitIgnoreGenerated := gitignore.GenerateGitIgnore()

	// Check the result of the .gitignore generation
	if gitIgnoreGenerated {
		log.Info(".gitignore successfully generated")
	} else {
		log.Fatalf("error generating .gitignore file")
	}

	// calling the function to execute Terraform command
	log.WithField("command", args.terraformCommand).Info("executing Terraform command")
	executeTerraformCommand(args.terraformCommand, accessKey, secretKey, sessionToken)

}

func getSSMParameters(client aws.SSMClient, paramKeys []string) map[string]string {
	params := make(map[string]string)

	for _, key := range paramKeys {
		param, err := aws.GetSSMParameter(client, key)
		if err != nil {
			handleError(err, "failed to get SSM Parameter for key: "+key)
			return nil
		}
		params[key] = param
	}

	return params
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
		os.Exit(1)
	}
}

func executeTerraformCommand(terraformCommand, accessKey, secretKey, sessionToken string) {
	commandWithArgs := strings.Fields(terraformCommand)
	terraformCmd := exec.Command("terraform", commandWithArgs...)
	terraformCmd.Env = append(os.Environ(),
		"AWS_ACCESS_KEY_ID="+accessKey,
		"AWS_SECRET_ACCESS_KEY="+secretKey,
		"AWS_SESSION_TOKEN="+sessionToken,
	)

	var stdout, stderr bytes.Buffer
	terraformCmd.Stdout = &stdout
	terraformCmd.Stderr = &stderr

	err := terraformCmd.Run()
	handleError(err, "cmd.Run() failed")

	if err != nil {
		log.Fatalf("cmd.Run() failed: %s\nStderr: %s", err, stderr.String())
	}

	fmt.Printf("output:\n%s\n", stdout.String())
}

// Define function to set tfS3Key based on the current directory
func getTFS3Key(pwd string, targetAccount string) (string, error) {
	var tfS3Key string

	// Check if the current directory is "aft-account-customizations" and set tfS3Key accordingly
	if strings.Contains(pwd, "aft-account-customizations") {
		tfS3Key = fmt.Sprintf("%s-aft-account-customizations/terraform.tfstate", targetAccount)
	} else if strings.Contains(pwd, "aft-global-customizations") {
		// Check if the current directory is "aft-global-customizations" and set tfS3Key accordingly
		tfS3Key = fmt.Sprintf("%s-aft-global-customizations/terraform.tfstate", targetAccount)
	} else { // If the directory is neither, return an error
		log.Errorf("Run aftctl local from the aft-account-customizations or aft-global-customizations repository.")
		return "", fmt.Errorf("invalid directory")
	}

	return tfS3Key, nil
}

func processJinjaFiles(
	tfDistributionParam string,
	ctMgmtRegionParam string,
	tfBackendRegionParam string,
	aftMgmtAccountIDParam string,
	aftExecutionRoleNameParam string,
	tfS3BucketIDParam string,
	tfS3Key string,
	tfDynamoDBTableNameParam string,
	tfKmsKeyIDParam string,
) error {

	// Read all files with the extension ".jinja"
	files, err := filepath.Glob("*.jinja")
	if err != nil {
		return fmt.Errorf("error reading jinja files: %v", err)
	}

	for _, f := range files {
		// Read the jinja template file
		input, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("error reading file:", f, err)
			continue
		}

		// Parse the jinja template
		template, err := pongo2.FromString(string(input))
		if err != nil {
			fmt.Println("error parsing template:", err)
			continue
		}

		// Define constants and generate timestamp
		partition := "aws"
		timestamp := time.Now().Format(time.RFC3339)

		// Execute the template with context
		output, err := template.Execute(pongo2.Context{
			"timestamp":             timestamp,
			"tf_distribution_type":  tfDistributionParam,
			"provider_region":       ctMgmtRegionParam,
			"region":                tfBackendRegionParam,
			"aft_admin_role_arn":    fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, aftMgmtAccountIDParam, aftExecutionRoleNameParam),
			"target_admin_role_arn": fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, args.targetAccount, aftExecutionRoleNameParam),
			"bucket":                tfS3BucketIDParam,
			"key":                   tfS3Key,
			"dynamodb_table":        tfDynamoDBTableNameParam,
			"kms_key_id":            tfKmsKeyIDParam,
		})
		if err != nil {
			fmt.Println("error executing template:", err)
			continue
		}

		// Write the output to a new file with the same name but different extension
		outputFile := fmt.Sprintf("./%s.tf", strings.TrimSuffix(f, ".jinja"))
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			fmt.Println("error writing output file:", err)
			continue
		}
	}

	return nil
}

func validateInput() error {
	_, err := validate.CheckAWSAccountID(args.targetAccount)
	if err != nil {
		return fmt.Errorf("invalid AWS Account ID: %w", err)
	}
	_, err = validate.CheckTerraformCommand(args.terraformCommand)
	if err != nil {
		return fmt.Errorf("invalid Terraform Command: %w", err)
	}
	return nil
}

func initializeAWSandSSMClients() (*aws.Client, aws.SSMClient, error) {
	log.Info("initializing AWS Client using AFT Account credentials")
	awsClient := aws.NewClient("")
	if awsClient == nil {
		return nil, nil, fmt.Errorf("failed to initialize AWS client")
	}

	ssmClient := awsClient.GetSSMClient()
	if ssmClient == nil {
		return nil, nil, fmt.Errorf("failed to initialize SSM client")
	}

	return awsClient, ssmClient, nil
}
