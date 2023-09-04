/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package version

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/edgarsilva948/aftctl/pkg/info"
)

// Cmd represents the Cobra command for the version functionality.
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the tool",
	Long:  "Prints the version number of the tool.",
	Run:   Run,
}

// Run executes the version command, printing the version of the tool.
func Run(cmd *cobra.Command, argv []string) {
	info.PrintVersion(os.Stdout)
}
