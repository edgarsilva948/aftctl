/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/edgarsilva948/aftctl/pkg/aws"

	"gopkg.in/yaml.v2"
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

	// aftctl deploy aft -f arquivo.yaml
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
		"The Name for metadata",
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

	cfg := aws.Config{
		Region: "us-east-1",
	}
	aws.InitAWSClient(cfg)

	s3Client := aws.NewS3Client()

	isValidName := regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString
	isValidBucketName := regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString

	// fileNameIsNull() {
	//	return args.filename != ""
	// }
	if args.filename != "" {

		// Allowed flags with -f
		allowedFlagsWithF := []string{"watch", "dry-run"}

		// Check for other flags
		flags := cmd.Flags()
		invalidFlagCombination := false
		flags.VisitAll(func(flag *pflag.Flag) {
			if flag.Name != "file" && flag.Changed {
				isAllowed := false
				for _, allowedFlag := range allowedFlagsWithF {
					if flag.Name == allowedFlag {
						isAllowed = true
						break
					}
				}

				if !isAllowed {
					fmt.Println("When -f flag is set, no other flags should be provided except for:", allowedFlagsWithF)
					invalidFlagCombination = true
				}
			}
		})

		if invalidFlagCombination {
			return
		}

		// reading the yaml file if informed with -f flag
		yamlFile, err := os.ReadFile(*&args.filename)
		if err != nil {
			log.Fatalf("Error reading YAML file: %s\n", err)
			return
		}

		// Unmarshal the YAML file
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("Error unmarshalling YAML: %s\n", err)
			return
		}

		// Name: not informed (file)
		if config.Metadata.Name == "" {
			config.Metadata.Name = args.name
		}

		// Name: informed correctly (file)
		if config.Metadata.Name != "" && isValidName(config.Metadata.Name) {
			fmt.Printf("Deploying the AFT using the Metadata name: %s\n", config.Metadata.Name)
		}

		// Name: informed incorrectly (file)
		if config.Metadata.Name != "" && !isValidName(config.Metadata.Name) {
			log.Fatalf("Metadata 'name' must be any combination of uppercase (A-Z), lowercase (a-z) alphabets, numbers (0-9) along with hyphens (-).")
			return
		}

		// Create Bucket: Option not informed (file)
		if config.DeploymentConfiguration.CreateTerraformStateBucket == nil {
			log.Fatalf("Deployment Configuration 'createTerraformStateBucket' is required")
			return
		}

		// Bucket Name: not informed (file)
		if config.DeploymentConfiguration.TerraformStateBucketName == "" {
			log.Fatalf("Deployment Configuration 'terraformStateBucketName' is required")
			return
		}

		// Bucket Name: informed incorrectly (file)
		if !isValidBucketName(config.DeploymentConfiguration.TerraformStateBucketName) {
			log.Fatalf("Deployment Configuration 'terraformStateBucketName' must be any combination of uppercase (A-Z), lowercase (a-z) alphabets, numbers (0-9) along with hyphens (-).")
			return
		}

		// Bucket Name: informed correctly (file) and the Create Bucket Option is false
		if isValidBucketName(config.DeploymentConfiguration.TerraformStateBucketName) && !*config.DeploymentConfiguration.CreateTerraformStateBucket {
			fmt.Printf("The bucket %s will be used to store tfstate file\n", config.DeploymentConfiguration.TerraformStateBucketName)
		}

		// Bucket Name: informed correctly (file) and the Create Bucket Option is true
		if isValidBucketName(config.DeploymentConfiguration.TerraformStateBucketName) && *config.DeploymentConfiguration.CreateTerraformStateBucket {
			CreateStateBucketInAftAccount(s3Client, config.DeploymentConfiguration.TerraformStateBucketName)
		}

	} else {

		// Name: informed correctly (flag)
		if args.name != "" && isValidName(args.name) {
			fmt.Printf("Deploying the AFT using the Metadata name: %s\n", args.name)
			config.Metadata.Name = args.name
		}

		// Name: informed incorrectly (flag)
		if args.name != "" && !isValidName(args.name) {
			log.Fatalf("flag '--name' (-n) must be any combination of uppercase (A-Z), lowercase (a-z) alphabets, numbers (0-9) along with hyphens (-).")
			return
		}

		// Create Bucket: Option informed (flag)
		if args.createTerraformStateBucket {
			*config.DeploymentConfiguration.CreateTerraformStateBucket = args.createTerraformStateBucket
		}

		// Bucket Name: not informed (flag)
		if args.terraformStateBucketName == "" {
			log.Fatalf("flag '--terraform-state-bucket-name' is required")
			return
		}

		// Bucket Name: informed incorrectly (flag)
		if !isValidBucketName(args.terraformStateBucketName) {
			log.Fatalf("flag '--terraform-state-bucket-name' must be any combination of uppercase (A-Z), lowercase (a-z) alphabets, numbers (0-9) along with hyphens (-).")
			return
		}

		// Bucket Name: informed correctly (flag) and the Create Bucket Option is false
		if isValidBucketName(args.terraformStateBucketName) && !args.createTerraformStateBucket {
			fmt.Printf("The bucket %s will be used to store tfstate file\n", args.terraformStateBucketName)
			config.DeploymentConfiguration.TerraformStateBucketName = args.terraformStateBucketName
		}

		// Bucket Name: informed correctly (flag) and the Create Bucket Option is true
		if isValidBucketName(args.terraformStateBucketName) && args.createTerraformStateBucket {
			config.DeploymentConfiguration.TerraformStateBucketName = args.terraformStateBucketName
			CreateStateBucketInAftAccount(s3Client, args.terraformStateBucketName)
		}

		//
		// needs to create all the validations for bucket path
		//
		config.DeploymentConfiguration.TerraformStateBucketPath = args.terraformStateBucketPath
	}
}

// CreateStateBucketInAftAccount creates a state bucket in the AFT account.
func CreateStateBucketInAftAccount(client aws.S3Client, bucketName string) (bool, error) {

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
	exists, err := aws.BucketExists(client, bucketName, "us-east-1")
	if err != nil {
		return false, fmt.Errorf("error checking if bucket exists: %w", err)
	}
	if exists {
		return false, fmt.Errorf("error: The bucket named '%s' already exists", bucketName)
	}
	return true, nil
}
