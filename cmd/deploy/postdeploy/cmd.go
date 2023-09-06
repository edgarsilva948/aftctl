/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package postdeploy

import "github.com/spf13/cobra"

// Cmd is the exported command for the AFT prerequisites.
var Cmd = &cobra.Command{
	Use:   "prereqs",
	Short: "Setup AFT post deployment in AFT-Management Account",
	Long:  "Setup AFT post deployment in AFT-Management Account",
	Example: `# aftctl usage examples"
	  aftctl deploy postdeploy -f deployment.yaml
	
	  aftctl deploy postdeploy --region="us-east-1"`,
	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false
}

func run(cmd *cobra.Command, _ []string) {
}
