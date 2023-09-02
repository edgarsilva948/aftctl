/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package version

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/edgarsilva948/aftctl/pkg/info"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the tool",
	Long:  "Prints the version number of the tool.",
	Run:   run,
}

func run(cmd *cobra.Command, argv []string) {
	fmt.Fprintf(os.Stdout, "%s\n", info.Version)
}
