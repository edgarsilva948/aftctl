/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// This file contains aws clients and session.
package aws

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client represents a client for Amazon S3.
type S3Client interface {
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

// A Config represents the AWS configuration settings.
type Config struct {
	Region string
}

// AWSSession is the AWS session used to interact with AWS services.
var AWSSession *session.Session

// InitializeAWS checks the AWS region and creates a new session.
func InitializeAWS() (bool, error) {
	isRegionSet, err := checkRegion()
	if !isRegionSet {
		log.Fatalf("Region Error: %s", err)
		return false, nil
	}

	isSessionSet, err := createSession(os.Getenv("AWS_REGION"))
	if !isSessionSet {
		log.Fatalf("AWS Session Error: %s", err)
		return false, nil
	}

	return true, nil
}

// checkRegion returns an error if region is not set
func checkRegion() (bool, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return false, errors.New("AWS_REGION environment variable not set")
	}
	return true, nil
}

// createSession returns a new session initialized
func createSession(region string) (bool, error) {
	var err error
	AWSSession, err = session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// NewS3Client returns a new S3Client initialized
func NewS3Client() *s3.S3 {
	return s3.New(AWSSession)
}
