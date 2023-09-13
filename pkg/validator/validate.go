/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package validate

import (
	"fmt"
	"strconv"
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
	var err error

	acceptedCommands := []string{"plan", "apply", "destroy", "init"}

	if command == "" {
		fmt.Printf("error: terraform command is required\n")
	}

	// Check if command is accepted
	for _, acceptedCommand := range acceptedCommands {
		if command == acceptedCommand {
			fmt.Println("command accepted")
			return true, nil
		} else {
			fmt.Printf("command is invalid: accepted commands are %s\n", acceptedCommands)
			return false, err
		}
	}

	return true, nil
}
