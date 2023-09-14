/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetSSMParameter from AFT Account
func GetSSMParameter(client SSMClient, paramName string) (string, error) {

	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	}

	result, err := client.GetParameter(input)

	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}
