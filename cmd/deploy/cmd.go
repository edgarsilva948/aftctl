/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package deploy

import (
	"github.com/edgarsilva948/aftctl/cmd/deploy/aft"
	"github.com/edgarsilva948/aftctl/pkg/arguments"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"setup"},
	Short:   "Deploy AFT from from stdin",
	Long:    "Deploy AFT from from stdin",
}

func init() {

	Cmd.AddCommand(aft.Cmd)

	flags := Cmd.PersistentFlags()
	arguments.AddProfileFlag(flags)
	arguments.AddRegionFlag(flags)

}
