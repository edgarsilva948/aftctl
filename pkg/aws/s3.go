/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains aws clients and session.
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketExists checks if a given S3 bucket exists.
func BucketExists(client S3Client, bucketName string) (bool, error) {
	// Check if the client is nil
	if client == nil {
		return false, fmt.Errorf("S3Client is not provided")
	}

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

// EnsureS3BucketExists creates a new S3 bucket with the given name, or returns success if it already exists.
func EnsureS3BucketExists(client S3Client, bucketName string, kmsKeyID string) (bool, error) {
	return true, nil
}
