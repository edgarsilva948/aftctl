package aws

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

type mockSessionCreator struct {
	shouldError bool
}

func (m *mockSessionCreator) CreateSession(region string) (*session.Session, error) {

	if region == "" {
		return nil, errors.New("region cannot be empty")
	}

	if m.shouldError {
		return nil, errors.New("mock error")
	}
	return &session.Session{}, nil
}

var _ = ginkgo.Describe("Interacting with AWS API", func() {

	ginkgo.Context("validating the checkRegion function", func() {

		ginkgo.When("AWS_REGION is not provided", func() {
			ginkgo.It("should return an error", func() {
				isSet, err := checkRegion()
				gomega.Expect(isSet).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("AWS_REGION environment variable not set"))
			})
		})

	})

	ginkgo.Context("validating the checkRegion function", func() {

		ginkgo.When("AWS_REGION is declared", func() {
			ginkgo.It("should return an error", func() {

				// Set the environment variable
				os.Setenv("AWS_REGION", "us-west-2")

				isSet, _ := checkRegion()
				gomega.Expect(isSet).To(gomega.BeTrue())
			})
		})

	})

	ginkgo.Context("Creating AWS session", func() {

		ginkgo.When("region is not informed", func() {
			ginkgo.It("should return nil and an error indicating the region is empty", func() {
				mock := &mockSessionCreator{shouldError: false}
				sess, err := mock.CreateSession("")
				gomega.Expect(sess).To(gomega.BeNil())
				gomega.Expect(err).ToNot(gomega.BeNil())
				gomega.Expect(err).To(gomega.MatchError("region cannot be empty"))
			})
		})

		ginkgo.When("region is valid and no errors occur", func() {
			ginkgo.It("should return a session and no error", func() {
				mock := &mockSessionCreator{shouldError: false}
				sess, err := mock.CreateSession("us-west-2")
				gomega.Expect(sess).ToNot(gomega.BeNil())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("an error occurs during session creation", func() {
			ginkgo.It("should return nil and an error", func() {
				mock := &mockSessionCreator{shouldError: true}
				sess, err := mock.CreateSession("us-west-2")
				gomega.Expect(sess).To(gomega.BeNil())
				gomega.Expect(err).ToNot(gomega.BeNil())
				gomega.Expect(err).To(gomega.MatchError("mock error"))
			})
		})

	})

})
