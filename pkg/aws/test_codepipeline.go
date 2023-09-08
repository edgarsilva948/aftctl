/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
package aws

import (
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/codepipeline/codepipelineiface"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// MockCodePipelineClient is a mock implementation of an S3 client for testing.
type MockCodePipelineClient struct {
	codepipelineiface.CodePipelineAPI
	CreatePipelineFunc func(*codepipeline.CreatePipelineInput) (*codepipeline.CreatePipelineOutput, error)
	ListPipelinesFunc  func(*codepipeline.ListPipelinesInput) (*codepipeline.ListPipelinesOutput, error)
}

// CreatePipeline is a mock implementation of the CreatePipeline method.
func (m *MockCodePipelineClient) CreatePipeline(input *codepipeline.CreatePipelineInput) (*codepipeline.CreatePipelineOutput, error) {
	return m.CreatePipelineFunc(input)
}

var _ = ginkgo.Describe("Interacting with the CodePipeline API", func() {

	ginkgo.Context("testing the checkIfCodePipelineClientIsProvided", func() {
		ginkgo.When("CodePipelineClient is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfCodePipelineClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("CodepipelineClient is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkIfPipelineNameIsProvided function", func() {
		ginkgo.When("pipelineName is not provided", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfPipelineNameIsProvided("")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("pipeline name is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkPipelineNameCompliance", func() {
		ginkgo.When("The Pipeline name is not between 3 and 100 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkPipelineNameCompliance("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("pipeline name must be between 3 and 100 characters long"))
			})
		})

		ginkgo.When("The pipeline name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkPipelineNameCompliance("pipeline!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("pipeline name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number"))
			})
		})

		ginkgo.When("The pipeline name is valid", func() {
			ginkgo.It("should return true", func() {
				isValid, err := checkPipelineNameCompliance("valid-pipeline-name")
				gomega.Expect(isValid).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

})
