/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package color

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var color string

var options = []string{"auto", "never", "always"}

// AddFlag adds the interactive flag to the given set of command line flags.
func AddFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&color,
		"color",
		"auto",
		fmt.Sprintf("Surround certain characters with escape sequences to display them in color "+
			"on the terminal. Allowed options are %s", options),
	)

	cmd.RegisterFlagCompletionFunc("color", completion)
}

func completion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return options, cobra.ShellCompDirectiveDefault
}

// UseColor returns a bool that indicates whether the color is enabled
func UseColor() bool {
	switch color {
	case "never":
		return false
	case "always":
		return true
	case "auto":
		fallthrough
	default:
		if runtime.GOOS == "windows" {
			return false
		}
		stdout, err := os.Stdout.Stat()
		if err != nil {
			return true
		}
		return (stdout.Mode()&os.ModeDevice != 0) && (stdout.Mode()&os.ModeNamedPipe == 0)
	}
}
