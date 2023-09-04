/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgarsilva948/aftctl/pkg/aws"
)

// Metadata holds the basic information.
type Metadata struct {
	Name string
}

// DeploymentConfiguration holds the terraform aditional resources
type DeploymentConfiguration struct {
	CreateTerraformStateBucket *bool  `yaml:"createTerraformStateBucket"`
	TerraformStateBucketName   string `yaml:"terraformStateBucketName"`
	TerraformStateBucketPath   string `yaml:"terraformStateBucketPath"`
}

// ControlTowerVariables holds the CT deployment information
type ControlTowerVariables struct {
	CTManagementAccountID    string `yaml:"controlTowerManagementAccountId"`
	LogArchiveAccountID      string `yaml:"logArchiveAccountId"`
	AuditAccountID           string `yaml:"auditAccountId"`
	AftManagementAccountID   string `yaml:"aftManagementAccountId"`
	CTHomeRegion             string `yaml:"controlTowerHomeRegion"`
	TFBackendSecondaryRegion string `yaml:"terraformBackendSecondaryRegion"`
}

// TerraformConfiguration holds the TF deployment information
type TerraformConfiguration struct {
	TerraformVersion      string `yaml:"terraformVersion"`
	TerraformDistribution string `yaml:"terraformDistribution"`
}

// VcsConfiguration holds the VCS deployment information
type VcsConfiguration struct {
	VcsProvider string `yaml:"vcsProvider"`
}

// Configuration holds the settings for the AFT deployment.
type Configuration struct {
	AftFeatureCloudtrailDataEvents     bool `yaml:"aftFeatureCloudtrailDataEvents"`
	AftFeatureEnterpriseSupport        bool `yaml:"aftFeatureEnterpriseSupport"`
	AftFeatureDeleteDefaultVpcsEnabled bool `yaml:"aftFeatureDeleteDefaultVpcsEnabled"`
}

// Config holds the full settings for the AFT deployment.
type Config struct {
	Metadata                Metadata
	DeploymentConfiguration DeploymentConfiguration `yaml:"deploymentConfiguration"`
	ControlTowerVariables   ControlTowerVariables   `yaml:"controlTowerVariables"`
	TerraformConfiguration  TerraformConfiguration  `yaml:"terraformConfiguration"`
	VcsConfiguration        VcsConfiguration        `yaml:"vcsConfiguration"`
	Configuration           Configuration           `yaml:"aftConfiguration"`
}

var args struct {
	// Watch logs during deployment
	watch bool

	// Simulate deploying AFT
	dryRun bool

	// The Deployment YAML file to process
	filename string

	// Metadata args
	name string

	// Deployment Configuration args
	terraformStateBucketName   string
	createTerraformStateBucket bool
	terraformStateBucketPath   string
}

// Cmd is the exported command for the AFT deployment.
var Cmd = &cobra.Command{
	Use:   "aft",
	Short: "Setup AFT in AFT-Management Account",
	Long:  "Setup AFT in AFT-Management Account",
	Example: `# aftctl usage examples"
	  aftctl deploy aft -f deployment.yaml
	
	  aftctl deploy aft --region="us-east-1"`,
	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false

	flags.BoolVarP(
		&args.watch,
		"watch",
		"w",
		false,
		"Watch logs during deployment.",
	)

	flags.BoolVar(
		&args.dryRun,
		"dry-run",
		false,
		"Simulate deploying AFT",
	)

	flags.StringVarP(
		&args.filename,
		"file",
		"f",
		"",
		"This file contains the deployment "+
			"instructions to deploy AFT",
	)

	flags.StringVarP(
		&args.name,
		"name",
		"n",
		"aft-deploy-configuration",
		"A metadata Name for the deployment",
	)

	flags.BoolVar(
		&args.createTerraformStateBucket,
		"create-terraform-state-bucket",
		true,
		"Whether to create a terraform state bucket",
	)

	flags.StringVar(
		&args.terraformStateBucketName,
		"terraform-state-bucket-name",
		"",
		"Name of the deployment terraform state bucket",
	)

	flags.StringVar(
		&args.terraformStateBucketPath,
		"terraform-state-bucket-path",
		"",
		"Path to save the state file inside the terraform state bucket",
	)
}

var config Config

func run(cmd *cobra.Command, _ []string) {

	s3Client := aws.NewS3Client()

	CreateStateBucketInAftAccount(s3Client, args.terraformStateBucketName)
}

// CreateStateBucketInAftAccount creates a state bucket in the AFT account.
func CreateStateBucketInAftAccount(client aws.S3Client, bucketName string) (bool, error) {

	isClientValid, err := isClientValid(client)
	if !isClientValid {
		fmt.Println("Raw error:", err)
		return false, err
	}

	isBucketNameValid, err := checkBucketName(bucketName)
	if !isBucketNameValid {
		return false, err
	}

	isBucketReady, err := checkBucketStatus(client, bucketName)
	if !isBucketReady {
		return false, err
	}

	return true, nil
}

func isClientValid(client aws.S3Client) (bool, error) {
	if client == nil {
		return false, errors.New("client is nil")
	}

	return true, nil
}

func checkBucketName(bucketName string) (bool, error) {
	length := len(bucketName)

	// Bucket names must be between 3 (min) and 63 (max) characters long.
	if length < 3 || length > 63 {
		return false, errors.New("bucket name must be between 3 and 63 characters long")
	}

	//Bucket names must not start with the prefix xn--.
	// Bucket names must not start with the prefix sthree- and the prefix sthree-configurator.
	if strings.HasPrefix(bucketName, "xn--") || strings.HasPrefix(bucketName, "sthree-") {
		return false, errors.New("bucket name cannot start with restricted prefixes (xn-- or sthree-)")
	}

	// Bucket names must not end with the suffix -s3alias. This suffix is reserved for access point alias names. For more information, see Using a bucket-style alias for your S3 bucket access point.
	// Bucket names must not end with the suffix --ol-s3. This suffix is reserved for Object Lambda Access Point alias names. For more information, see How to use a bucket-style alias for your S3 bucket Object Lambda Access Point.
	if strings.HasSuffix(bucketName, "-s3alias") || strings.HasSuffix(bucketName, "--ol-s3") {
		return false, errors.New("bucket name cannot end with restricted suffixes (-s3alias or --ol-s3)")
	}

	// Bucket names can consist only of lowercase letters, numbers, and hyphens (-).
	pattern := `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(bucketName) {
		return false, errors.New("bucket name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number")
	}

	// Additional check to make sure bucket names don't have two adjacent periods.
	if strings.Contains(bucketName, "..") {
		return false, errors.New("bucket name must not contain two adjacent periods")
	}

	// Check for IP address format (which is not allowed)
	ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`
	ipRe := regexp.MustCompile(ipPattern)
	if ipRe.MatchString(bucketName) {
		return false, errors.New("bucket name must not be formatted as an IP address")
	}

	return true, nil
}

func checkBucketStatus(client aws.S3Client, bucketName string) (bool, error) {
	exists, err := aws.BucketExists(client, bucketName)
	if err != nil {
		return false, fmt.Errorf("error checking if bucket exists: %w", err)
	}
	if exists {
		return false, fmt.Errorf("error: The bucket named '%s' already exists", bucketName)
	}
	return true, nil
}
