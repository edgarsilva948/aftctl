/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"github.com/edgarsilva948/aftctl/cmd/aft/deploy"
	"github.com/spf13/cobra"
)

// Cmd represents the root command for the "deploy" functionality.
var Cmd = &cobra.Command{
	Use:   "aft",
	Short: "Deploy AFT from from stdin",
	Long:  "Deploy AFT from from stdin",
}

func init() {

	Cmd.AddCommand(deploy.Cmd)
}
