/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// MockCodeBuildClient is a mock implementation of an S3 client for testing.
type MockCodeBuildClient struct {
	codebuildiface.CodeBuildAPI

	CreateProjectFunc func(*codebuild.CreateProjectInput) (*codebuild.CreateProjectOutput, error)
	ListProjectsFunc  func(*codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error)
}

// ListProjects is a mock implementation of the ListProjects method.
func (m *MockCodeBuildClient) ListProjects(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
	return m.ListProjectsFunc(input)
}

// CreateProject is a mock implementation of the CreateProject method.
func (m *MockCodeBuildClient) CreateProject(input *codebuild.CreateProjectInput) (*codebuild.CreateProjectOutput, error) {
	return m.CreateProjectFunc(input)
}

var _ = ginkgo.Describe("Interacting with the CodeBuild API", func() {

	ginkgo.Context("testing the EnsureCodeBuildProjectExists function", func() {

		ginkgo.When("Project already exists", func() {
			ginkgo.It("should return an success", func() {

				mockClient := &MockCodeBuildClient{
					CreateProjectFunc: func(input *codebuild.CreateProjectInput) (*codebuild.CreateProjectOutput, error) {
						return &codebuild.CreateProjectOutput{}, nil
					},
					ListProjectsFunc: func(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
						return &codebuild.ListProjectsOutput{}, nil
					},
				}

				ensure, err := EnsureCodeBuildProjectExists(mockClient, "000000000000", "test-docker-image", "test-project", "test-repo", "test-branch", "test-role")

				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("project doesn't exists", func() {
			ginkgo.It("should create the project", func() {

				mockClient := &MockCodeBuildClient{
					CreateProjectFunc: func(input *codebuild.CreateProjectInput) (*codebuild.CreateProjectOutput, error) {
						return &codebuild.CreateProjectOutput{}, nil
					},
					ListProjectsFunc: func(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
						return &codebuild.ListProjectsOutput{}, nil
					},
				}

				ensure, err := EnsureCodeBuildProjectExists(mockClient, "000000000000", "test-docker-image", "test-project", "test-repo", "test-branch", "test-role")

				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("testing the checkIfCodeBuildClientIsProvided", func() {
		ginkgo.When("CodeBuildClient is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfCodeBuildClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("CodeBuildClient is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkIfProjectNameIsProvided function", func() {
		ginkgo.When("projectName is not provided", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfProjectNameIsProvided("")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("project name is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkProjectNameCompliance", func() {
		ginkgo.When("The Project name is not between 3 and 255 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkProjectNameCompliance("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("project name must be between 3 and 255 characters long"))
			})
		})

		ginkgo.When("The Project name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkProjectNameCompliance("project!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("project name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number"))
			})
		})

		ginkgo.When("The project name is valid", func() {
			ginkgo.It("should return true", func() {
				isValid, err := checkProjectNameCompliance("valid-project-name")
				gomega.Expect(isValid).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("testing the checkIfProjectExists function", func() {
		ginkgo.When("if the project already exists", func() {
			ginkgo.It("should return an error", func() {

				mockClient := &MockCodeBuildClient{
					ListProjectsFunc: func(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
						return &codebuild.ListProjectsOutput{
							Projects: []*string{
								aws.String("one-project"),
							},
						}, nil
					},
				}

				ensure, err := checkIfProjectExists(mockClient, "one-project")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("list operation fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockCodeBuildClient{
					ListProjectsFunc: func(input *codebuild.ListProjectsInput) (*codebuild.ListProjectsOutput, error) {
						return &codebuild.ListProjectsOutput{
							Projects: []*string{
								aws.String("any-project"),
							},
						}, nil
					},
				}

				exists, err := checkIfProjectExists(mockClient, "one-project")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
