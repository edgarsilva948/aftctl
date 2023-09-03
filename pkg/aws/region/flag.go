// This file contains functions used to implement the '--region' command line option.

package region

import (
	"os"

	"github.com/edgarsilva948/aftctl/pkg/helper"
	"github.com/spf13/pflag"
)

// AddFlag adds the debug flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&region,
		"region",
		"",
		"Use a specific AWS region, overriding the AWS_REGION environment variable.",
	)
}

// Region returns a string with the name of the AWS region being used.
func Region() string {
	if helper.HandleEscapedEmptyString(region) != "" {
		return region
	}
	awsRegion := os.Getenv("AWS_REGION")
	if helper.HandleEscapedEmptyString(awsRegion) != "" {
		return awsRegion
	}
	return ""
}

// region is a string flag that indicates which AWS region is being used.
var region string
