package aws

import (
	"errors"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// MockS3Client is a mock implementation of an S3 client for testing.
type MockS3Client struct {
	ListBucketsFunc func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

// ListBuckets is a mock implementation of the ListBuckets method.
func (m *MockS3Client) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(input)
}

var _ = ginkgo.Describe("Interacting with the S3 API", func() {

	ginkgo.Context("BucketExists function", func() {

		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				exists, err := BucketExists(nil, "bucketName")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("S3Client is not provided"))
			})
		})

		ginkgo.When("S3 list bucket operation fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return nil, errors.New("AWS S3 error")
					},
				}
				exists, err := BucketExists(mockClient, "bucketName")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("failed to list S3 buckets: AWS S3 error"))
			})
		})

		ginkgo.When("Bucket exists", func() {
			ginkgo.It("should return true", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("bucketName")},
							},
						}, nil
					},
				}
				exists, err := BucketExists(mockClient, "bucketName")
				gomega.Expect(exists).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("Bucket does not exist", func() {
			ginkgo.It("should return false", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("anotherBucket")},
							},
						}, nil
					},
				}
				exists, err := BucketExists(mockClient, "bucketName")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
