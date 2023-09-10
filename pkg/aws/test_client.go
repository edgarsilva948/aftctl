/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
	"github.com/aws/aws-sdk-go/service/codepipeline/codepipelineiface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// ClientMockS3Client is a mock of S3API
type ClientMockS3Client struct {
	s3iface.S3API
}

// ClientMockIamClient is a mock of IAMAPI
type ClientMockIamClient struct {
	iamiface.IAMAPI
}

// ClientMockCodePipelineClient is a mock of CodePipelineAPI
type ClientMockCodePipelineClient struct {
	codepipelineiface.CodePipelineAPI
}

// ClientMockCodeCommitClient is a mock of CodeCommitAPI
type ClientMockCodeCommitClient struct {
	codecommitiface.CodeCommitAPI
}

// ClientMockCodeBuildClient is a mock of CodeBuildAPI
type ClientMockCodeBuildClient struct {
	codebuildiface.CodeBuildAPI
}

// ClientMockCloudFormationClient is a mock of CloudFormationAPI
type ClientMockCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

var _ = ginkgo.Describe("Interacting with AWS API", func() {

	// Local variables for mock clients
	var (
		mockS3Client             *ClientMockS3Client
		mockIamClient            *ClientMockIamClient
		mockCodePipelineClient   *ClientMockCodePipelineClient
		mockCodeCommitClient     *ClientMockCodeCommitClient
		mockCodeBuildClient      *ClientMockCodeBuildClient
		mockCloudFormationClient *ClientMockCloudFormationClient
		client                   *Client
	)

	// BeforeEach setup for test
	ginkgo.BeforeEach(func() {
		mockS3Client = &ClientMockS3Client{}
		mockIamClient = &ClientMockIamClient{}
		mockCodePipelineClient = &ClientMockCodePipelineClient{}
		mockCodeCommitClient = &ClientMockCodeCommitClient{}
		mockCodeBuildClient = &ClientMockCodeBuildClient{}
		mockCloudFormationClient = &ClientMockCloudFormationClient{}

		// Initialize client with mock clients
		client = &Client{
			s3Client:             mockS3Client,
			iamClient:            mockIamClient,
			codepipelineClient:   mockCodePipelineClient,
			codecommitClient:     mockCodeCommitClient,
			codebuildClient:      mockCodeBuildClient,
			cloudformationClient: mockCloudFormationClient,
		}
	})

	// Context for fetching AWS Clients
	ginkgo.Context("Fetching AWS Clients", func() {

		ginkgo.When("GetS3Client is called", func() {
			ginkgo.It("should return the S3 client", func() {
				gomega.Expect(client.GetS3Client()).To(gomega.Equal(mockS3Client))
			})
		})

		ginkgo.When("GetIamClient is called", func() {
			ginkgo.It("should return the IAM client", func() {
				gomega.Expect(client.GetIamClient()).To(gomega.Equal(mockIamClient))
			})
		})

		ginkgo.When("GetCodePipelineClient is called", func() {
			ginkgo.It("should return the CodePipeline client", func() {
				gomega.Expect(client.GetCodePipelineClient()).To(gomega.Equal(mockCodePipelineClient))
			})
		})

		ginkgo.When("GetCodeCommitClient is called", func() {
			ginkgo.It("should return the CodeCommit client", func() {
				gomega.Expect(client.GetCodeCommitClient()).To(gomega.Equal(mockCodeCommitClient))
			})
		})

		ginkgo.When("GetCodeBuildClient is called", func() {
			ginkgo.It("should return the CodeBuild client", func() {
				gomega.Expect(client.GetCodeBuildClient()).To(gomega.Equal(mockCodeBuildClient))
			})
		})

		ginkgo.When("GetCloudFormationClient is called", func() {
			ginkgo.It("should return the CloudFormation client", func() {
				gomega.Expect(client.GetCloudFormationClient()).To(gomega.Equal(mockCloudFormationClient))
			})
		})
	})
})
