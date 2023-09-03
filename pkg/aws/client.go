/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aws

import (
	"log"

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

// InitAWSClient initializes the AWS client
func InitAWSClient(cfg Config) {
	var err error
	AWSSession, err = session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region)},
	)

	if err != nil {
		log.Fatalf("AWS Session Error: %s", err)
	}
}

// NewS3Client returns a new S3Client initialized
func NewS3Client() *s3.S3 {
	return s3.New(AWSSession)
}
