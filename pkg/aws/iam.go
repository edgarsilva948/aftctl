/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aws

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/edgarsilva948/aftctl/pkg/aws/tags"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// EnsureIamRoleExists creates a new IAM Role with the given name, or returns success if it already exists.
func EnsureIamRoleExists(client IAMClient, roleName string, trustRelationShipService string, policyName string, region string, aftAccount string, repoName string, bucketName string, terraformStateBucketName string) (bool, error) {

	_, err := checkIfIamClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfRoleNameIsProvided(roleName)

	if err != nil {
		return false, err
	}

	roleExists, _ := checkIfRoleExists(client, roleName)

	if !roleExists {
		fmt.Printf("IAM Role %s doesn't exists... creating\n", roleName)

		_, err := createRole(
			client,
			roleName,
			trustRelationShipService,
			policyName,
			region,
			aftAccount,
			repoName,
			bucketName,
			terraformStateBucketName,
		)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true
	logger, _ := config.Build()

	defer logger.Sync()

	message := fmt.Sprintf("IAM Role %s already exists", roleName)

	customIamInfoLog(logger, message)

	return true, nil
}

func customIamInfoLog(logger *zap.Logger, msg string) {

	// Security-related emoji
	lockEmoji := "🔒"
	coloredMsg := "\x1b[32m" + lockEmoji + " " + msg + "\x1b[0m"

	logger.Info(coloredMsg)
}

// func to verify if the given iam role is provided
func checkIfRoleNameIsProvided(roleName string) (bool, error) {
	if roleName == "" {
		fmt.Printf("Error: %v\n", "role name is not provided")
		return false, fmt.Errorf("role name is not provided")
	}

	isRoleNameValid, err := checkRoleNameCompliance(roleName)
	if !isRoleNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given iam is compliant
func checkRoleNameCompliance(roleName string) (bool, error) {
	length := len(roleName)
	// iam names must be between 3 (min) and 63 (max) characters long.
	if length < 3 || length > 64 {
		return false, errors.New("iam name must be between 3 and 63 characters long")
	}

	// iam names can consist only of lowercase letters, numbers, and hyphens (-).
	pattern := `^[\w+=,.@-]{1,64}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(roleName) {
		return false, errors.New("iam name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number")
	}

	return true, nil
}

// func to verify if the given role name already exists
func checkIfRoleExists(client IAMClient, roleName string) (bool, error) {

	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	_, err := client.GetRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false, nil
			default:
				return false, aerr
			}
		}
		return false, err
	}
	return true, nil
}

// func to verify if the given client is valid
func checkIfIamClientIsProvided(client IAMClient) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("IAMClient is not provided")
	}

	return true, nil
}

// func to create given role if it doesn't exist'
func createRole(client IAMClient, roleName string, trustRelationShipService string, policyName string, region string, aftAccount string, repoName string, bucketName string, terraformStateBucketName string) (bool, error) {

	// Assume Role Policy Document
	const iamAssumeRolePolicyDocument = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "%s"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`

	// Role Policy
	const rolePolicyDocument = `{
		"Version":"2012-10-17",
		"Statement":[
		   {
			  "Resource":"*",
			  "Effect":"Allow",
			  "Action":[
				 "codebuild:StartBuild",
				 "codebuild:BatchGetBuilds"
			  ]
		   },
		   {
			"Resource":"*",
			"Effect":"Allow",
			"Action":[
			   "logs:CreateLogGroup",
			   "logs:CreateLogStream",
			   "logs:PutLogEvents"
			]
 		   },
		   {
			  "Resource":"arn:aws:codecommit:%s:%s:%s",
			  "Effect":"Allow",
			  "Action":[
				 "codecommit:GetBranch",
				 "codecommit:GetCommit",
				 "codecommit:UploadArchive",
				 "codecommit:GetUploadArchiveStatus",
				 "codecommit:CancelUploadArchive"
			  ]
		   },
		   {
			"Effect": "Allow",
			"Resource": "arn:aws:s3:::%s/*",
			"Action": [
				"s3:PutObject",
				"s3:GetObject",
				"s3:GetObjectVersion",
				"s3:GetBucketVersioning"
			]
			},
			{
				"Effect": "Allow",
				"Resource": "arn:aws:s3:::%s/*",
				"Action": [
					"s3:PutObject",
					"s3:GetObject",
					"s3:GetObjectVersion",
					"s3:GetBucketVersioning"
				]
			},
		   {
			  "Resource":"*",
			  "Effect":"Allow",
			  "Action":[
				 "ec2:CreateNetworkInterface",
				 "ec2:DescribeDhcpOptions",
				 "ec2:DescribeNetworkInterfaces",
				 "ec2:DeleteNetworkInterface",
				 "ec2:DescribeSubnets",
				 "ec2:DescribeSecurityGroups",
				 "ec2:DescribeVpcs",
				 "ec2:CreateNetworkInterfacePermission"
			  ]
		   }
		]
	 }`

	createRoleInput := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(fmt.Sprintf(iamAssumeRolePolicyDocument, trustRelationShipService)),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(roleName),
		Tags: []*iam.Tag{
			{
				Key:   aws.String(tags.Aftctl),
				Value: aws.String(tags.True),
			},
		},
	}

	_, err := client.CreateRole(createRoleInput)
	if err != nil {
		log.Printf("unable to create role %q, %v", roleName, err)
		return false, err
	}

	putPolicyInput := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(fmt.Sprintf(rolePolicyDocument, region, aftAccount, repoName, bucketName, terraformStateBucketName)),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	}

	_, err = client.PutRolePolicy(putPolicyInput)
	if err != nil {
		return false, err
	}

	fmt.Printf("IAM Role %s successfully created\n", roleName)
	return true, nil
}
