/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CheckAWSAccountID checks if a string represents a valid AWS account id
func CheckAWSAccountID(accountID string) (bool, error) {
	var err error
	if len(accountID) != 12 {
		fmt.Printf("error: account id must be 12 characters long\n")
		return false, err
	}

	// Check if all characters are digits
	_, err = strconv.ParseUint(accountID, 10, 64)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CheckTerraformCommand checks if a string represents a valid AWS account id
func CheckTerraformCommand(command string) (bool, error) {

	if command == "" {
		return false, errors.New("error: terraform command is required")
	}

	acceptedCommands := []string{
		"plan",
		"apply",
		"apply --auto-approve",
		"destroy",
		"init",
		"import",
		"validate",
		"output",
		"fmt",
		"force-unlock",
		"get",
		"providers",
		"refresh",
		"show",
		"state",
		"taint",
		"untaint",
		"version",
		"workspace",
	}

	for _, acceptedCommand := range acceptedCommands {
		// Check if the beginning of the command string matches an accepted command
		if strings.HasPrefix(command, acceptedCommand) {
			return true, nil
		}
	}

	fmt.Printf("command is invalid: accepted commands are %s\n", acceptedCommands)
	return false, errors.New("invalid command")

}
