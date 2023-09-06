/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"github.com/spf13/cobra"
)

// Metadata holds the basic information.
type Metadata struct {
	Name string
}

// DeploymentConfiguration holds the terraform aditional resources
type DeploymentConfiguration struct {
	TerraformStateBucketName string `yaml:"terraformStateBucketName"`
	TerraformStateBucketPath string `yaml:"terraformStateBucketPath"`
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
}

var config Config

func run(cmd *cobra.Command, _ []string) {

}
