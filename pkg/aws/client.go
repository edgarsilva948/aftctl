/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains aws clients and session.
package aws

import (
	"fmt"
	"os"

	"github.com/caarlos0/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/codepipeline/codepipelineiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// S3Client represents a client for Amazon S3.
type S3Client interface {
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	WaitUntilBucketExists(*s3.HeadBucketInput) error
	PutPublicAccessBlock(*s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error)
	PutBucketPolicy(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)
	PutBucketTagging(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// CodeCommitClient represents a client for Amazon Code Commit.
type CodeCommitClient interface {
	CreateRepository(*codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error)
	GetRepository(*codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error)
	TagResource(*codecommit.TagResourceInput) (*codecommit.TagResourceOutput, error)
}

// CodeBuildClient represents a client for Amazon Code Build.
type CodeBuildClient interface {
	CreateProject(*codebuild.CreateProjectInput) (*codebuild.CreateProjectOutput, error)
	ListProjects(*codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error)
}

// IAMClient represents a client for Amazon Code Commit.
type IAMClient interface {
	CreateRole(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	PutRolePolicy(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	GetRole(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
}

// CodePipelineClient represents a client for Amazon Code Pipeline.
type CodePipelineClient interface {
	CreatePipeline(*codepipeline.CreatePipelineInput) (*codepipeline.CreatePipelineOutput, error)
	ListPipelines(input *codepipeline.ListPipelinesInput) (*codepipeline.ListPipelinesOutput, error)
}

// CloudformationClient represents a client for Cloudformation.
type CloudformationClient interface {
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

// SSMClient represents a client for SSM.
type SSMClient interface {
	GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

// STSClient represents a client for STS.
type STSClient interface {
	AssumeRole(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
}

// Client struct implementing all the client interfaces
type Client struct {
	s3Client             s3iface.S3API
	iamClient            iamiface.IAMAPI
	codepipelineClient   codepipelineiface.CodePipelineAPI
	codecommitClient     codecommitiface.CodeCommitAPI
	codebuildClient      codebuildiface.CodeBuildAPI
	cloudformationClient cloudformationiface.CloudFormationAPI
	ssmClient            ssmiface.SSMAPI
	stsClient            stsiface.STSAPI
}

// NewClient loads credentials following the chain credentials
func NewClient(profile string) *Client {

	opts := session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}

	sess, err := session.NewSessionWithOptions(opts)

	// Check for session initialization error
	if err != nil {
		fmt.Println("error creating session:", err)
		os.Exit(1)
	}

	// Check for a nil session
	if sess == nil {
		log.Fatal("session is nil")
		os.Exit(1)
	}

	// Check for a nil Config
	if sess.Config == nil {
		log.Fatal("invalid session configuration")
		os.Exit(1)
	}

	// Check for nil Credentials
	if sess.Config.Credentials == nil {
		log.Fatal("invalid session credentials")
		os.Exit(1)
	}

	// Check for credential errors
	_, errCreds := sess.Config.Credentials.Get()
	if errCreds != nil {
		fmt.Println("credential error:", errCreds)
		os.Exit(1)
	}

	// Check for an unset AWS region
	if aws.StringValue(sess.Config.Region) == "" {
		fmt.Println("region is not set.")
		os.Exit(1)
	}

	// Check for an unset aws default profile and a profile set in the environment variable
	if awsProfile := os.Getenv("AWS_PROFILE"); awsProfile != "" && profile == "" {
		log.WithField("profile", awsProfile).Info("using AWS_PROFILE environment variable")
	}

	return &Client{
		s3Client:             s3.New(sess),
		iamClient:            iam.New(sess),
		codepipelineClient:   codepipeline.New(sess),
		codecommitClient:     codecommit.New(sess),
		codebuildClient:      codebuild.New(sess),
		cloudformationClient: cloudformation.New(sess),
		ssmClient:            ssm.New(sess),
		stsClient:            sts.New(sess),
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

// GetCodeBuildClient returns the client for AWS CodeBuild service.
func (ac *Client) GetCodeBuildClient() codebuildiface.CodeBuildAPI {
	return ac.codebuildClient
}

// GetCloudFormationClient returns the client for AWS CloudFormation service.
func (ac *Client) GetCloudFormationClient() cloudformationiface.CloudFormationAPI {
	return ac.cloudformationClient
}

// GetSSMClient returns the client for AWS SSM service.
func (ac *Client) GetSSMClient() ssmiface.SSMAPI {
	return ac.ssmClient
}

// GetSTSClient returns the client for AWS STS service.
func (ac *Client) GetSTSClient() stsiface.STSAPI {
	return ac.stsClient
}

// GetAWSCredentials returns the AWS credentials for the given profile.
func GetAWSCredentials(profile string) (string, string, string, error) {

	opts := session.Options{
		Profile: profile,
	}

	sess, err := session.NewSessionWithOptions(opts)

	if err != nil {
		return "", "", "", err
	}

	creds := sess.Config.Credentials
	credValue, err := creds.Get()
	if err != nil {
		return "", "", "", err
	}

	return credValue.AccessKeyID, credValue.SecretAccessKey, credValue.SessionToken, nil
}
