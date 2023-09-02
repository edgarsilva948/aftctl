/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains functions that add common arguments to the command line.

package arguments

import (
	"github.com/edgarsilva948/aftctl/pkg/debug"
	"github.com/spf13/pflag"
)

// AddDebugFlag adds the '--debug' flag to the given set of command line flags.
func AddDebugFlag(fs *pflag.FlagSet) {
	debug.AddFlag(fs)
}
