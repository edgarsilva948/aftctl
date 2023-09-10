/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edgarsilva948/aftctl/cmd/aft"
	"github.com/edgarsilva948/aftctl/cmd/completion"
	"github.com/edgarsilva948/aftctl/cmd/docs"
	"github.com/edgarsilva948/aftctl/cmd/version"

	"github.com/edgarsilva948/aftctl/pkg/color"
)

var root = &cobra.Command{
	Use:   "aftctl",
	Short: "Command line tool for Amazon Account Factory for Terraform (AFT).",
	Long: "Command line tool for Amazon Account Factory for Terraform (AFT).\n" +
		"For further documentation visit " +
		"https://docs.aws.amazon.com/controltower/latest/userguide/aft-getting-started.html\n",
}

func init() {
	// Add the command line flags:
	color.AddFlag(root)

	// Register the subcommands:
	root.AddCommand(completion.Cmd)
	root.AddCommand(docs.Cmd)
	root.AddCommand(version.Cmd)
	root.AddCommand(aft.Cmd)
}

func main() {
	// Execute the root command:
	root.SetArgs(os.Args[1:])
	err := root.Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "Did you mean this?") {
			fmt.Fprintf(os.Stderr, "Failed to execute root command: %s\n", err)
		}
		os.Exit(1)
	}
}
