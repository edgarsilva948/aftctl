/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package local provides the local command
package local

import (
	"errors"
	"os"

	awsAft "github.com/edgarsilva948/aftctl/pkg/aws"
	profile "github.com/edgarsilva948/aftctl/pkg/aws/profiles"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// MockSTSClient is a mock implementation of an STS client for testing.
type MockSTSClient struct {
	stsiface.STSAPI
	AssumeRoleFunc func(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
}

// MockSSMClient is a mock implementation of an SSM client for testing.
type MockSSMClient struct {
	ssmiface.SSMAPI
	GetParameterFunc func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

// GetParameter is a mock implementation of the GetParameter method.
func (m *MockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return m.GetParameterFunc(input)
}

// AssumeRole is a mock implementation of the AssumeRole method.
func (m *MockSTSClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return m.AssumeRoleFunc(input)
}

var _ = ginkgo.Describe("Interacting with the SSM API", func() {
	var (
		originalWd string
		err        error
	)

	// Use BeforeEach and AfterEach for setup and teardown specific to this Describe block
	ginkgo.BeforeEach(func() {
		originalWd, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		err := os.Chdir(originalWd)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.Context("testing the GetSSMParameter function", func() {

		var mockClient *MockSSMClient

		ginkgo.BeforeEach(func() {
			mockClient = &MockSSMClient{}
		})

		ginkgo.When("parameter retrieval is successful", func() {
			ginkgo.It("should return the parameter value", func() {
				mockClient.GetParameterFunc = func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
					return &ssm.GetParameterOutput{
						Parameter: &ssm.Parameter{
							Value: aws.String("some_value"),
						},
					}, nil
				}

				value, err := awsAft.GetSSMParameter(mockClient, "some_parameter")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(value).To(gomega.Equal("some_value"))
			})
		})

		ginkgo.When("parameter retrieval fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient.GetParameterFunc = func(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
					return nil, errors.New("some error")
				}
				value, err := awsAft.GetSSMParameter(mockClient, "some_parameter")
				gomega.Expect(err).ToNot(gomega.BeNil())
				gomega.Expect(err.Error()).To(gomega.Equal("some error"))
				gomega.Expect(value).To(gomega.BeEmpty())
			})
		})

	})

	ginkgo.Context("when getting the current directory", func() {
		ginkgo.When("there is no error", func() {
			ginkgo.It("should return the current directory", func() {
				pwd, err := os.Getwd()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(pwd).To(gomega.Equal(originalWd))
			})
		})

		ginkgo.When("there is an error", func() {
			ginkgo.It("should handle the error", func() {

				err := os.Chdir("/non/existent/directory")
				gomega.Expect(err).To(gomega.HaveOccurred())

				_, err = os.Getwd()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

	})

	ginkgo.Context("when setting up the profile", func() {
		var mockSTSClient *MockSTSClient
		var aftMgmtAccountIDParam string
		var aftAdminRoleNameParam string

		ginkgo.BeforeEach(func() {
			mockSTSClient = &MockSTSClient{
				AssumeRoleFunc: func(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
					return &sts.AssumeRoleOutput{}, nil
				},
			}
			aftMgmtAccountIDParam = "someAccountID"
			aftAdminRoleNameParam = "someRoleName"
		})

		ginkgo.When("there is no error", func() {
			ginkgo.It("should complete without error", func() {

				mockSTSClient = &MockSTSClient{
					AssumeRoleFunc: func(*sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
						return &sts.AssumeRoleOutput{
							Credentials: &sts.Credentials{
								AccessKeyId:     aws.String("some-access-key-id"),
								SecretAccessKey: aws.String("some-secret-access-key"),
								SessionToken:    aws.String("some-session-token"),
							},
						}, nil
					},
				}
				err := profile.SetupProfile(mockSTSClient, aftMgmtAccountIDParam, aftAdminRoleNameParam, "AWSAFT-Session")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.When("there is an error", func() {
			ginkgo.BeforeEach(func() {
				mockSTSClient.AssumeRoleFunc = func(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
					return nil, errors.New("some error")
				}
			})

			ginkgo.It("should handle the error", func() {
				err := profile.SetupProfile(mockSTSClient, aftMgmtAccountIDParam, aftAdminRoleNameParam, "AWSAFT-Session")
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.Equal("some error"))
			})
		})
	})

})
