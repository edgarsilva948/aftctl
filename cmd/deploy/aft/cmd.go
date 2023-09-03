/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/edgarsilva948/aftctl/pkg/aws"

	"gopkg.in/yaml.v2"
)

// Regions supported by AWS Control Tower
const (
	// ilCentral1 represents Israel (Tel Aviv)
	ilCentral1 = "il-central-1"

	// afSouth1 represents Africa (Cape Town)
	afSouth1 = "af-south-1"

	// apEast1 represents Asia Pacific (Hong Kong)
	apEast1 = "ap-east-1"

	// apNortheast3 represents Asia Pacific (Osaka)
	apNortheast3 = "ap-northeast-3"

	// apSoutheast3 represents Asia Pacific (Jakarta)
	apSoutheast3 = "ap-southeast-3"

	// euSouth1 represents Europe (Milan)
	euSouth1 = "eu-south-1"

	// meSouth1 represents Middle East (Bahrain)
	meSouth1 = "me-south-1"

	// usWest1 represents US West (N. California)
	usWest1 = "us-west-1"

	// usGovEast1 represents AWS GovCloud (US-East)
	usGovEast1 = "us-gov-east-1"

	// usGovWest1 represents AWS GovCloud (US-West)
	usGovWest1 = "us-gov-west-1"

	// euWest3 represents Europe (Paris)
	euWest3 = "eu-west-3"

	// saEast1 represents South America (Sao Paulo)
	saEast1 = "sa-east-1"

	// apNortheast1 represents Asia Pacific (Tokyo)
	apNortheast1 = "ap-northeast-1"

	// apNortheast2 represents Asia Pacific (Seoul)
	apNortheast2 = "ap-northeast-2"

	// apSouth1 represents Asia Pacific (Mumbai)
	apSouth1 = "ap-south-1"

	// apSoutheast1 represents Asia Pacific (Singapore)
	apSoutheast1 = "ap-southeast-1"

	// caCentral1 represents Canada (Central)
	caCentral1 = "ca-central-1"

	// euCentral1 represents Europe (Frankfurt)
	euCentral1 = "eu-central-1"

	// euNorth1 represents Europe (Stockholm)
	euNorth1 = "eu-north-1"

	// euWest2 represents Europe (London)
	euWest2 = "eu-west-2"

	// apSoutheast2 represents Asia Pacific (Sydney)
	apSoutheast2 = "ap-southeast-2"

	// usEast1 represents US East (N. Virginia)
	usEast1 = "us-east-1"

	// usWest2 represents US West (Oregon)
	usWest2 = "us-west-2"

	// euWest1 represents Europe (Ireland)
	euWest1 = "eu-west-1"

	// usEast2 represents US East (Ohio)
	usEast2 = "us-east-2"
)

type Metadata struct {
	Name string
}

type DeploymentConfiguration struct {
	CreateTerraformStateBucket *bool  `yaml:"createTerraformStateBucket"`
	TerraformStateBucketName   string `yaml:"terraformStateBucketName"`
	TerraformStateBucketPath   string `yaml:"terraformStateBucketPath"`
}

type ControlTowerVariables struct {
	CTManagementAccountID    string `yaml:"controlTowerManagementAccountId"`
	LogArchiveAccountID      string `yaml:"logArchiveAccountId"`
	AuditAccountID           string `yaml:"auditAccountId"`
	AftManagementAccountID   string `yaml:"aftManagementAccountId"`
	CTHomeRegion             string `yaml:"controlTowerHomeRegion"`
	TFBackendSecondaryRegion string `yaml:"terraformBackendSecondaryRegion"`
}

type TerraformConfiguration struct {
	TerraformVersion      string `yaml:"terraformVersion"`
	TerraformDistribution string `yaml:"terraformDistribution"`
}

type VcsConfiguration struct {
	VcsProvider string `yaml:"vcsProvider"`
}

type AftConfiguration struct {
	AftFeatureCloudtrailDataEvents     bool `yaml:"aftFeatureCloudtrailDataEvents"`
	AftFeatureEnterpriseSupport        bool `yaml:"aftFeatureEnterpriseSupport"`
	AftFeatureDeleteDefaultVpcsEnabled bool `yaml:"aftFeatureDeleteDefaultVpcsEnabled"`
}

type Config struct {
	Metadata                Metadata
	DeploymentConfiguration DeploymentConfiguration `yaml:"deploymentConfiguration"`
	ControlTowerVariables   ControlTowerVariables   `yaml:"controlTowerVariables"`
	TerraformConfiguration  TerraformConfiguration  `yaml:"terraformConfiguration"`
	VcsConfiguration        VcsConfiguration        `yaml:"vcsConfiguration"`
	AftConfiguration        AftConfiguration        `yaml:"aftConfiguration"`
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

var validRegions = []string{
	ilCentral1,
	afSouth1,
	apEast1,
	apNortheast3,
	apSoutheast3,
	euSouth1,
	meSouth1,
	usWest1,
	usGovEast1,
	usGovWest1,
	euWest3,
	saEast1,
	apNortheast1,
	apNortheast2,
	apSouth1,
	apSoutheast1,
	caCentral1,
	euCentral1,
	euNorth1,
	euWest2,
	apSoutheast2,
	usEast1,
	usWest2,
	euWest1,
	usEast2,
}

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

	// flags.StringVar(
	// 	&args.controlTowerManagementAccountId,
	// 	"controltower-management-account-id",
	// 	"",
	// 	"The Management Account ID that will be used during the deployment process.",
	// )
}

var config Config

func run(cmd *cobra.Command, _ []string) {

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
			createStateBucketInAftAccount(config.DeploymentConfiguration.TerraformStateBucketName)
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
			createStateBucketInAftAccount(args.terraformStateBucketName)
		}

		//
		// needs to create all the validations for bucket path
		//
		config.DeploymentConfiguration.TerraformStateBucketPath = args.terraformStateBucketPath
	}
}

// Functions
func createStateBucketInAftAccount(bucketName string) {

	exists, err := aws.BucketExists(bucketName, "us-east-1")
	if err != nil {
		fmt.Println("Error checking if bucket exists:", err)
		return
	}

	if exists {
		log.Fatalf("Error: The bucket named '%s' already exists. To use an existing bucket, please set the 'createTerraformStateBucket' or 'create-terraform-state-bucket' configuration to 'false'.", bucketName)
	} else {

		exampleKms := ""

		fmt.Printf("Bucket %s does not exist, creating", bucketName)
		err := aws.CreateS3Bucket(bucketName, exampleKms, "us-east-1")
		if err != nil {
			fmt.Printf("Error while creating the bucket: %s\n", err)
		} else {
			fmt.Println("Successfully created the bucket")
		}
	}
}
