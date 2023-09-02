/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains functions used to implement the '--profile' command line option.

package profile

import (
	"os"

	"github.com/spf13/pflag"
)

// AddFlag adds the debug flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&profile,
		"profile",
		"",
		"Use a specific AWS profile from your credential file.",
	)
}

// Profile returns a string with the name of the AWS profile being used.
func Profile() string {
	if profile != "" {
		return profile
	}
	awsProfile := os.Getenv("AWS_PROFILE")
	if awsProfile != "" {
		return awsProfile
	}
	return ""
}

// profile is a string flag that indicates which AWS profile is being used.
var profile string
