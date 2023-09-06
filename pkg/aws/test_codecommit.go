/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
package aws

import (
	"errors"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/codecommit/codecommitiface"
)

// MockCodeCommitClient is a mock implementation of an S3 client for testing.
type MockCodeCommitClient struct {
	codecommitiface.CodeCommitAPI

	CreateRepositoryFunc func(*codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error)
	GetRepositoryFunc    func(*codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error)
}

// GetRepository is a mock implementation of the GetRepository method.
func (m *MockCodeCommitClient) GetRepository(input *codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error) {
	return m.GetRepositoryFunc(input)
}

// CreateRepository is a mock implementation of the CreateRepository method.
func (m *MockCodeCommitClient) CreateRepository(input *codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error) {
	return m.CreateRepositoryFunc(input)
}

var _ = ginkgo.Describe("Interacting with the CodeCommit API", func() {

	ginkgo.Context("testing the EnsureCodeCommitRepoExists function", func() {

		ginkgo.When("repository already exists", func() {
			ginkgo.It("should return a success", func() {

				mockClient := &MockCodeCommitClient{
					GetRepositoryFunc: func(input *codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error) {
						return &codecommit.GetRepositoryOutput{
							RepositoryMetadata: &codecommit.RepositoryMetadata{
								RepositoryName: input.RepositoryName,
							},
						}, nil
					},
				}

				ensure, err := EnsureCodeCommitRepoExists(mockClient, "repo", "simple description")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("repository doesn't exist", func() {
			ginkgo.It("should create the repository", func() {

				mockClient := &MockCodeCommitClient{
					CreateRepositoryFunc: func(input *codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error) {
						return &codecommit.CreateRepositoryOutput{}, nil
					},
					GetRepositoryFunc: func(input *codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error) {
						return nil, awserr.New(codecommit.ErrCodeRepositoryDoesNotExistException, "Repository does not exist", nil)
					},
				}

				ensure, err := EnsureCodeCommitRepoExists(mockClient, "repo", "simple description")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("repository creation fails", func() {
			ginkgo.It("should return an error", func() {

				mockClient := &MockCodeCommitClient{
					CreateRepositoryFunc: func(input *codecommit.CreateRepositoryInput) (*codecommit.CreateRepositoryOutput, error) {
						return nil, errors.New("AWS create repository error")
					},
					GetRepositoryFunc: func(input *codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error) {
						return nil, awserr.New(codecommit.ErrCodeRepositoryDoesNotExistException, "Repository does not exist", nil)
					},
				}

				ensure, err := EnsureCodeCommitRepoExists(mockClient, "failed-repo", "simple description")
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("AWS create repository error"))
			})
		})

		ginkgo.Context("testing the repoExists function", func() {

			ginkgo.When("CodeCommit list repo operation fails", func() {
				ginkgo.It("should return an error", func() {
					mockClient := &MockCodeCommitClient{
						GetRepositoryFunc: func(input *codecommit.GetRepositoryInput) (*codecommit.GetRepositoryOutput, error) {
							return nil, awserr.New(codecommit.ErrCodeRepositoryDoesNotExistException, "failed to list CodeCommit repositories", nil)
						},
					}
					exists, err := repoExists(mockClient, "repo")
					gomega.Expect(exists).To(gomega.BeFalse())
					gomega.Expect(err).To(gomega.MatchError("RepositoryDoesNotExistException: failed to list CodeCommit repositories"))
				})
			})

		})

	})

	ginkgo.Context("testing the checkIfCodeCommitClientIsProvided", func() {
		ginkgo.When("CodeCommitClient is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfCodeCommitClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("CodeCommitClient is not provided"))
			})
		})
	})

})
