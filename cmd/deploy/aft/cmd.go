/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// Regions
const (
	// RegionUSWest1 represents the US West Region North California
	RegionUSWest1 = "us-west-1"
)

type Config struct {
	Metadata                Metadata
	DeploymentConfiguration DeploymentConfiguration `yaml:"deploymentConfiguration"`
	ControlTowerVariables   ControlTowerVariables   `yaml:"controlTowerVariables"`
	TerraformConfiguration  TerraformConfiguration  `yaml:"terraformConfiguration"`
	VcsConfiguration        VcsConfiguration        `yaml:"vcsConfiguration"`
	AftConfiguration        AftConfiguration        `yaml:"aftConfiguration"`
}

type Metadata struct {
	Name   string
	Region string
}

type DeploymentConfiguration struct {
	CreateTerraformStateBucket bool   `yaml:"createTerraformStateBucket"`
	TFStateBucketName          string `yaml:"terraformStateBucketName"`
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

var args struct {
	// Watch logs during deployment
	watch bool

	// Simulate deploying AFT
	dryRun bool

	// The Deployment YAML file to process
	filename string

	// CT Management account ID
	controlTowerManagementAccountId string

	// CT Log Archive account ID
	logArchiveAccountId string

	// CT Audit account ID
	auditAccountId string

	// AFT account ID
	aftManagementAccountId string

	// Control Tower main region
	controlTowerHomeRegion string

	// Control Tower seccondary region
	terraformBackendSecondaryRegion string
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

	flags.StringVarP(
		&args.filename,
		"file",
		"f",
		"",
		"This file contains the deployment"+
			"instructions to deploy AFT",
	)

	flags.BoolVar(
		&args.dryRun,
		"dry-run",
		false,
		"Simulate deploying AFT",
	)

	flags.BoolVarP(
		&args.watch,
		"watch",
		"w",
		false,
		"Watch cluster installation logs.",
	)

	flags.StringVar(
		&args.controlTowerManagementAccountId,
		"controltower-management-account-id",
		"",
		"The Management Account ID that will be used during the deployment process.",
	)
}

func run(cmd *cobra.Command, _ []string) {

	// Allowed flags with -f
	allowedFlagsWithF := []string{"watch", "dry-run"}

	// Check if the -f flag is set
	if &args.filename == nil {
		fmt.Println("Please provide a YAML file using -f flag")
		return
	}

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

	yamlFile, err := os.ReadFile(*&args.filename)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
		return
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
	}

	fmt.Printf("Successfully read file: %s\n", *&args.filename)
}
