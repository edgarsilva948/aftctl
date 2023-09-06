/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains aws clients and session.
package aws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/codepipeline/codepipelineiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// S3Client represents a client for Amazon S3.
type S3Client interface {
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	WaitUntilBucketExists(*s3.HeadBucketInput) error
	PutPublicAccessBlock(*s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error)
	PutBucketPolicy(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)
	PutBucketTagging(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)
}

// CodeCommitClient represents a client for Amazon Code Commit.
type CodeCommitClient interface {
	CreateRepository(*codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error)
	GetRepository(*codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error)
	TagResource(*codecommit.TagResourceInput) (*codecommit.TagResourceOutput, error)
}

// IAMClient represents a client for Amazon Code Commit.
type IAMClient interface {
	CreateRole(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	PutRolePolicy(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
}

// Client struct implementing all the client interfaces
type Client struct {
	s3Client           s3iface.S3API
	iamClient          iamiface.IAMAPI
	codepipelineClient codepipelineiface.CodePipelineAPI
	codecommitClient   codecommitiface.CodeCommitAPI
}

// NewClient loads credentials following the chain credentials
func NewClient() *Client {

	sess, err := session.NewSession()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if sess == nil {
		log.Fatal("Session is nil")
	}

	_, errCreds := sess.Config.Credentials.Get()
	if errCreds != nil {
		fmt.Println(errCreds)
		os.Exit(1)
	}

	if aws.StringValue(sess.Config.Region) == "" {
		fmt.Println("Region is not set.")
		os.Exit(1)
	}

	return &Client{
		s3Client:           s3.New(sess),
		iamClient:          iam.New(sess),
		codepipelineClient: codepipeline.New(sess),
		codecommitClient:   codecommit.New(sess),
	}
}

// GetS3Client fetches the S3 Client and enables the cmd to use
func (ac *Client) GetS3Client() s3iface.S3API {
	return ac.s3Client
}

// GetIamClient fetches the IAM Client and enables the cmd to use
func (ac *Client) GetIamClient() iamiface.IAMAPI {
	return ac.iamClient
}

// GetCodePipelineClient fetches the Code Pipeline Client and enables the cmd to use
func (ac *Client) GetCodePipelineClient() codepipelineiface.CodePipelineAPI {
	return ac.codepipelineClient
}

// GetCodeCommitClient fetches the CodeCommit Client and enables the cmd to use
func (ac *Client) GetCodeCommitClient() codecommitiface.CodeCommitAPI {
	return ac.codecommitClient
}
