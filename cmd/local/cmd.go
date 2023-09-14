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

	"log"
	"os"

	"path/filepath"
	"strings"

	"github.com/edgarsilva948/aftctl/pkg/aws"
	profile "github.com/edgarsilva948/aftctl/pkg/aws/profiles"
	"github.com/edgarsilva948/aftctl/pkg/gitignore"
	"github.com/edgarsilva948/aftctl/pkg/logging"
	validate "github.com/edgarsilva948/aftctl/pkg/validator"
	"github.com/flosch/pongo2"

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

const assumeRoleIcon = "🔄"
const profileIcon = "📝"
const gitIgnoreIcon = "⚙️"
const terraformIcon = "⛏️ "

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
	// defining the S3 key for the local execution
	var tfS3Key string

	// validate input account ID
	_, err := validate.CheckAWSAccountID(args.targetAccount)
	if err != nil {
		os.Exit(1)
	}

	// validate input terraform command
	_, err = validate.CheckTerraformCommand(args.terraformCommand)
	if err != nil {
		os.Exit(1)
	}

	// client initialization with AFT Credentials
	message := "Initializing AWS Client using AFT Account credentials... step (1/4)"
	logging.CustomLog(profileIcon, "green", message)

	awsClient := aws.NewClient("")

	aftMgmtAccountIDParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), aftMgmtAccountID)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	tfDistributionParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), tfDistribution)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	ctMgmtRegionParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), ctMgmtRegion)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	aftAdminRoleNameParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), aftAdminRoleName)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}
	aftExecutionRoleNameParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), aftExecutionRoleName)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	tfBackendRegionParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), tfBackendRegion)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	tfKmsKeyIDParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), tfKmsKeyID)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	tfDynamoDBTableNameParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), tfDynamoDBTableName)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	tfS3BucketIDParam, err := aws.GetSSMParameter(awsClient.GetSSMClient(), tfS3BucketID)
	if err != nil {
		log.Fatalf("Failed to get SSM Parameter: %v", err)
		os.Exit(1)
	}

	// check if the current directory is aft-account-customizations or aft-global-customizations
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// check if the current directory is aft-account-customizations or aft-global-customizations and defining the S3 key for the local execution
	if strings.Contains(pwd, "aft-account-customizations") {
		tfS3Key = fmt.Sprintf("%s-aft-account-customizations/terraform.tfstate", args.targetAccount)
	} else if strings.Contains(pwd, "aft-global-customizations") {
		tfS3Key = fmt.Sprintf("%s-aft-global-customizations/terraform.tfstate", args.targetAccount)
	} else {
		fmt.Println("Run aftctl local from the aft-account-customizations or aft-global-customizations repository.")
		return
	}

	// setting up the AWS profile for the AFT Account using the user current credentials
	if err := profile.SetupProfile(awsClient.GetSTSClient(), aftMgmtAccountIDParam, aftAdminRoleNameParam, "AWSAFT-Session"); err != nil {
		fmt.Println("Error setting up profile", err)
	}

	message = fmt.Sprintf("Successfully set up AWS profile %s-%s Step (2/4)", aftMgmtAccountIDParam, aftAdminRoleNameParam)
	logging.CustomLog(assumeRoleIcon, "green", message)

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
		log.Fatalf("Error setting env var AWS_PROFILE: %v", err)
	}

	files, err := filepath.Glob("*.jinja")
	if err != nil {
		fmt.Println("Error reading jinja files:", err)
		return
	}

	for _, f := range files {
		input, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("Error reading file:", f, err)
			continue
		}

		template, err := pongo2.FromString(string(input))
		if err != nil {
			fmt.Println("Error parsing template:", err)
			continue
		}

		partition := "aws"
		timestamp := time.Now().Format(time.RFC3339)

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
			fmt.Println("Error executing template:", err)
			continue
		}

		outputFile := fmt.Sprintf("./%s.tf", strings.TrimSuffix(f, ".jinja"))
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			fmt.Println("Error writing output file:", err)
			continue
		}
	}

	gitIgnoreGenerated := gitignore.GenerateGitIgnore()

	if !gitIgnoreGenerated {
		log.Fatalf("error generating .gitignore file")
	}

	if gitIgnoreGenerated {
		message = ".gitignore successfully generated... (3/4)"
		logging.CustomLog(gitIgnoreIcon, "green", message)
	}

	message = fmt.Sprintf("Executing Terraform command %s... (4/4)", args.terraformCommand)
	logging.CustomLog(terraformIcon, "green", message)

	commandWithArgs := strings.Fields(args.terraformCommand)
	terraformCmd := exec.Command("terraform", commandWithArgs...)
	terraformCmd.Env = append(os.Environ(),
		"AWS_ACCESS_KEY_ID="+accessKey,
		"AWS_SECRET_ACCESS_KEY="+secretKey,
		"AWS_SESSION_TOKEN="+sessionToken,
	)

	var stdout, stderr bytes.Buffer
	terraformCmd.Stdout = &stdout
	terraformCmd.Stderr = &stderr

	err = terraformCmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed: %s\nStderr: %s", err, stderr.String())
	}

	fmt.Printf("Output:\n%s\n", stdout.String())

}
