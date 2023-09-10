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
	"github.com/edgarsilva948/aftctl/pkg/logging"
)

// FileReader is a function type that reads a file and returns its content as a byte slice.
type FileReader func(string) ([]byte, error)

// ReadFile is a variable of type FileReader, initially set to os.ReadFile.
var ReadFile FileReader = os.ReadFile

const uploadIcon = "⬆️ "

// UploadToS3 to upload the zip file to S3 Bucket
func UploadToS3(client S3Client, bucketName string, bucketKey string, fileName string) error {
	// Read the file
	fileContent, err := ReadFile(fileName)
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

	message := fmt.Sprintf("zip file %s successfully uploaded", fileName)
	logging.CustomLog(uploadIcon, "green", message)

	return nil
}
