/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains aws clients and session.
package aws

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/edgarsilva948/aftctl/pkg/aws/tags"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// EnsureS3BucketExists creates a new S3 bucket with the given name, or returns success if it already exists.
func EnsureS3BucketExists(client S3Client, bucketName string, aftManagementAccountID string, kmsKeyID string) (bool, error) {

	_, err := checkIfS3ClientIsProvided(client)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = checkIfBucketNameIsProvided(bucketName)

	if err != nil {
		return false, err
	}

	bucketExists, _ := bucketExists(client, bucketName)

	if !bucketExists {
		fmt.Printf("S3 bucket %s doesn't exists... creating\n", bucketName)

		_, err := createBucket(client, bucketName, aftManagementAccountID, kmsKeyID)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	fmt.Printf("S3 bucket %s already exists... continuing\n", bucketName)

	return true, nil
}

// BucketExists checks if a given S3 bucket exists.
func bucketExists(client S3Client, bucketName string) (bool, error) {

	isBucketExistent, err := checkIfBucketExists(client, bucketName)
	if err != nil {
		return false, err
	}

	return isBucketExistent, nil
}

// func to verify if the given client is valid
func checkIfS3ClientIsProvided(client S3Client) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("S3Client is not provided")
	}

	return true, nil
}

// func to verify if the given bucket name already exists
func checkIfBucketExists(client S3Client, bucketName string) (bool, error) {
	input := &s3.ListBucketsInput{}

	output, err := client.ListBuckets(input)
	if err != nil {
		return false, fmt.Errorf("failed to list S3 buckets: %w", err)
	}

	for _, bucket := range output.Buckets {
		if *bucket.Name == bucketName {
			return true, nil
		}
	}

	return false, nil
}

// func to verify if the given bucket is provided
func checkIfBucketNameIsProvided(bucketName string) (bool, error) {
	if bucketName == "" {
		fmt.Printf("Error: %v\n", "bucket name is not provided")
		return false, fmt.Errorf("bucket name is not provided")
	}

	isBucketNameValid, err := checkBucketNameCompliance(bucketName)
	if !isBucketNameValid {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	return true, nil
}

// func to verify if the given bucket is compliant
func checkBucketNameCompliance(bucketName string) (bool, error) {
	length := len(bucketName)

	// Bucket names must be between 3 (min) and 63 (max) characters long.
	if length < 3 || length > 63 {
		return false, errors.New("bucket name must be between 3 and 63 characters long")
	}

	//Bucket names must not start with the prefix xn--.
	// Bucket names must not start with the prefix sthree- and the prefix sthree-configurator.
	if strings.HasPrefix(bucketName, "xn--") || strings.HasPrefix(bucketName, "sthree-") {
		return false, errors.New("bucket name cannot start with restricted prefixes (xn-- or sthree-)")
	}

	// Bucket names must not end with the suffix -s3alias. This suffix is reserved for access point alias names. For more information, see Using a bucket-style alias for your S3 bucket access point.
	// Bucket names must not end with the suffix --ol-s3. This suffix is reserved for Object Lambda Access Point alias names. For more information, see How to use a bucket-style alias for your S3 bucket Object Lambda Access Point.
	if strings.HasSuffix(bucketName, "-s3alias") || strings.HasSuffix(bucketName, "--ol-s3") {
		return false, errors.New("bucket name cannot end with restricted suffixes (-s3alias or --ol-s3)")
	}

	// Bucket names can consist only of lowercase letters, numbers, and hyphens (-).
	pattern := `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(bucketName) {
		return false, errors.New("bucket name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number")
	}

	// Additional check to make sure bucket names don't have two adjacent periods.
	if strings.Contains(bucketName, "..") {
		return false, errors.New("bucket name must not contain two adjacent periods")
	}

	// Check for IP address format (which is not allowed)
	ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`
	ipRe := regexp.MustCompile(ipPattern)
	if ipRe.MatchString(bucketName) {
		return false, errors.New("bucket name must not be formatted as an IP address")
	}

	return true, nil
}

// func to create given bucket if it doesn't exist'
func createBucket(client S3Client, bucketName string, aftManagementAccountID string, kmsKeyID string) (bool, error) {

	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Printf("unable to create bucket %q, %v", bucketName, err)
		return false, err
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Printf("error occurred while waiting for bucket to be created, %v: %v", bucketName, err)
		return false, err
	}

	_, err = client.PutPublicAccessBlock(&s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			IgnorePublicAcls:      aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(true),
			RestrictPublicBuckets: aws.Bool(true),
		},
	})

	if err != nil {
		return false, err
	}
	_, err = client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(fmt.Sprintf(WriteAndListPolicyTemplateForAccount, aftManagementAccountID, bucketName, bucketName)),
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	_, err = client.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String(tags.Aftctl),
					Value: aws.String(tags.True),
				},
			},
		},
	})
	if err != nil {
		return false, err
	}

	fmt.Printf("S3 Bucket %s successfuly created", bucketName)
	return true, nil
}
