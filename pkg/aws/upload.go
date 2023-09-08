/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains aws clients and session.
package aws

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadToS3 to upload the zip file to S3 Bucket
func UploadToS3(client S3Client, bucketName string, bucketKey string, fileName string) error {
	// Read the file
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file %w", err)
	}

	// Upload the file
	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
		Body:   bytes.NewReader(fileContent),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %w", err)
	}

	return nil
}
