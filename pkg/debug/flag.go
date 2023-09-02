/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains functions used to implement the '--debug' command line option.

package debug

import (
	"github.com/spf13/pflag"
)

// AddFlag adds the debug flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVar(
		&enabled,
		"debug",
		false,
		"Enable debug mode.",
	)
}

// Enabled returns a boolean flag that indicates if the debug mode is enabled.
func Enabled() bool {
	return enabled
}

// enabled is a boolean flag that indicates that the debug mode is enabled.
var enabled bool
