/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package local

import (
	"os"

	validate "github.com/edgarsilva948/aftctl/pkg/validator"
	"github.com/spf13/cobra"
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
	Short: "Runs local AFT execution",
	Long:  "Runs AFT locally executing the same commands as the pipeline",
	Run:   Run,
}

// Run executes the local command
func Run(cmd *cobra.Command, argv []string) {
	_, err := validate.CheckAWSAccountID(args.targetAccount)
	if err != nil {
		os.Exit(1)
	}

	_, err = validate.CheckTerraformCommand(args.terraformCommand)
	if err != nil {
		os.Exit(1)
	}

}
