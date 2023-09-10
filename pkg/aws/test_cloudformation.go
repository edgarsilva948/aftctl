/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
package aws

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

// MockCloudformationClient is a mock implementation of an Cloudformation for testing.
type MockCloudformationClient struct {
	cloudformationiface.CloudFormationAPI

	CreateStackFunc    func(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DescribeStacksFunc func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

// DescribeStacks is a mock implementation of the DescribeStacks method.
func (m *MockCloudformationClient) DescribeStacks(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	return m.DescribeStacksFunc(input)
}

// CreateStack is a mock implementation of the CreateStack method.
func (m *MockCloudformationClient) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return m.CreateStackFunc(input)
}

var _ = ginkgo.Describe("Interacting with the Cloudformation API", func() {

	ginkgo.Context("testing the EnsureCloudformationExists function", func() {

		ginkgo.When("stack already exists", func() {
			ginkgo.It("should return an success", func() {

				mockClient := &MockCloudformationClient{
					DescribeStacksFunc: func(input *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
						return &cloudformation.DescribeStacksOutput{
							Stacks: []*cloudformation.Stack{
								{
									StackName: aws.String("test-stack"),
								},
							},
						}, nil
					},
				}

				ensure, err := EnsureCloudformationExists(mockClient, "test-stack", "test-repo-name", "test-description", "test-stack-name", "test-zip-filename")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("testing the checkIfCloudformationClientIsProvided", func() {
		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfCloudformationClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("CloudformationClient is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkIfStackNameIsProvided function", func() {
		ginkgo.When("stackName is not provided", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfStackNameIsProvided("")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("stack name is not provided"))
			})
		})

		ginkgo.When("The stack name is not between 3 and 100 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkStackNameCompliance("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("stack name must be between 3 and 100 characters long"))
			})
		})

		ginkgo.When("The Stack name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkStackNameCompliance("stack!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("stack name can only consist of lowercase letters, numbers, and hyphens"))
			})
		})

		ginkgo.When("The Stack name is valid", func() {
			ginkgo.It("should return true", func() {
				isValid, err := checkStackNameCompliance("valid-stack-name")
				gomega.Expect(isValid).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

})
