// This file contains functions that add common arguments to the command line.

package arguments

import (
	"github.com/edgarsilva948/aftctl/pkg/aws/profile"
	"github.com/edgarsilva948/aftctl/pkg/aws/region"
	"github.com/edgarsilva948/aftctl/pkg/debug"

	"github.com/spf13/pflag"
)

// AddProfileFlag adds the '--profile' flag to the given set of command line flags.
func AddProfileFlag(fs *pflag.FlagSet) {
	profile.AddFlag(fs)
}

// AddRegionFlag adds the '--region' flag to the given set of command line flags.
func AddRegionFlag(fs *pflag.FlagSet) {
	region.AddFlag(fs)
}

// AddDebugFlag adds the '--debug' flag to the given set of command line flags.
func AddDebugFlag(fs *pflag.FlagSet) {
	debug.AddFlag(fs)
}
