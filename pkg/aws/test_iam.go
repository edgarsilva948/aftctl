/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for AWS clients and session.
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// MockIAMClient is a mock implementation of an IAM client for testing.
type MockIAMClient struct {
	iamiface.IAMAPI
	CreateRoleFunc    func(input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	PutRolePolicyFunc func(input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	GetRoleFunc       func(input *iam.GetRoleInput) (*iam.GetRoleOutput, error)
}

// CreateRole is a mock implementation of the CreateRole method.
func (m *MockIAMClient) CreateRole(input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
	if m.CreateRoleFunc != nil {
		return m.CreateRoleFunc(input)
	}
	return nil, nil
}

// PutRolePolicy is a mock implementation of the PutRolePolicy method.
func (m *MockIAMClient) PutRolePolicy(input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
	if m.PutRolePolicyFunc != nil {
		return m.PutRolePolicyFunc(input)
	}
	return nil, nil
}

// GetRole is a mock implementation of the GetRoleFunc method.
func (m *MockIAMClient) GetRole(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
	if m.PutRolePolicyFunc != nil {
		return m.GetRoleFunc(input)
	}
	return nil, nil
}

var _ = ginkgo.Describe("Interacting with the IAM API", func() {

	ginkgo.Context("testing the EnsureIamRoleExists function", func() {

		ginkgo.When("role already exists", func() {
			ginkgo.It("should return success without errors", func() {

				mockClient := &MockIAMClient{
					GetRoleFunc: func(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
						return &iam.GetRoleOutput{
							Role: &iam.Role{
								RoleName: aws.String("roleName"),
							},
						}, nil
					},
				}

				roleExists, err := EnsureIamRoleExists(mockClient, "test-role", "test-policy", "test-bucket", "test-input", "test-input", "test-input", "test-input", "test-bucket")

				gomega.Expect(roleExists).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())

			})
		})

		ginkgo.When("IAM client is not provided", func() {
			ginkgo.It("should return an error", func() {

				roleExists, err := EnsureIamRoleExists(nil, "test-role", "test-policy", "test-bucket", "test-input", "test-input", "test-input", "test-input", "test-bucket")

				gomega.Expect(roleExists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("IAMClient is not provided"))

			})
		})

		ginkgo.When("IAM Role Name is not provided", func() {
			ginkgo.It("should return an error", func() {

				mockClient := &MockIAMClient{
					GetRoleFunc: func(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
						return &iam.GetRoleOutput{
							Role: &iam.Role{
								RoleName: aws.String("roleName"),
							},
						}, nil
					},
				}

				roleExists, err := EnsureIamRoleExists(mockClient, "", "test-policy", "test-bucket", "test-input", "test-input", "test-input", "test-input", "test-bucket")

				gomega.Expect(roleExists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("role name is not provided"))

			})
		})

		ginkgo.When("Bucket doesn't exists", func() {
			ginkgo.It("should create the bucket", func() {

				mockClient := &MockIAMClient{
					GetRoleFunc: func(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
						return &iam.GetRoleOutput{
							Role: &iam.Role{
								RoleName: aws.String(""),
							},
						}, nil
					},
				}

				roleExists, err := EnsureIamRoleExists(mockClient, "new-role", "test-policy", "test-bucket", "test-input", "test-input", "test-input", "test-input", "test-bucket")

				gomega.Expect(roleExists).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())

			})
		})
	})

	ginkgo.Context("testing the checkIfRoleNameIsProvided function", func() {
		ginkgo.When("roleName is not provided", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfRoleNameIsProvided("")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("role name is not provided"))
			})
		})

		ginkgo.When(" RoleName is provided and not Valid ", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfRoleNameIsProvided("ab")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("iam name must be between 3 and 63 characters long"))
			})
		})
	})

	ginkgo.Context("testing the checkRoleNameCompliance", func() {
		ginkgo.When("The iam name is not between 3 and 63 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkRoleNameCompliance("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("iam name must be between 3 and 63 characters long"))
			})
		})

		ginkgo.When("The iam name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkRoleNameCompliance("iam!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("iam name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number"))
			})
		})

	})

	ginkgo.Context("testing the checkIfRoleExists function", func() {
		ginkgo.When("the role already exists", func() {
			ginkgo.It("should not return an error", func() {

				mockClient := &MockIAMClient{

					GetRoleFunc: func(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
						return &iam.GetRoleOutput{
							Role: &iam.Role{
								RoleName: aws.String("role-name"),
								Arn:      aws.String("arn:aws:iam::account-ID-without-hyphens:role/role-name"),
							},
						}, nil
					},
				}

				roleExists, err := checkIfRoleExists(mockClient, "role-name")

				gomega.Expect(roleExists).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Context("testing the checkIfIamClientIsProvided function", func() {
		ginkgo.Context("testing the checkIfIamClientIsProvided", func() {
			ginkgo.When("IamClient is not provided", func() {
				ginkgo.It("should return an error", func() {
					ensure, err := checkIfIamClientIsProvided(nil)
					gomega.Expect(ensure).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.MatchError("IAMClient is not provided"))
				})
			})
		})

	})

})
