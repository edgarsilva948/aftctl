/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

// Package aws contains tests for aws clients and session.
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

	ginkgo.Context("testing the EnsureS3bucketExists function", func() {

		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				creates, err := EnsureS3BucketExists(nil, "bucketName", "test-kms-key-id")
				gomega.Expect(creates).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("S3Client is not provided"))
			})
		})

		ginkgo.When("bucketName is not provided", func() {
			ginkgo.It("should return an error", func() {

				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("anotherBucket")},
							},
						}, nil
					},
				}
				creates, err := EnsureS3BucketExists(mockClient, "", "test-kms-key-id")
				gomega.Expect(creates).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name is not provided"))
			})
		})

	})

	ginkgo.Context("testing the bucketExists function", func() {

		ginkgo.When("bucketName is not provided", func() {
			ginkgo.It("should return an error", func() {

				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("anotherBucket")},
							},
						}, nil
					},
				}
				exists, err := bucketExists(mockClient, "")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name is not provided"))
			})
		})

		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return false and an error", func() {
				exists, err := bucketExists(nil, "bucketName")
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
				exists, err := bucketExists(mockClient, "bucketname")
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
								{Name: aws.String("bucketname")},
							},
						}, nil
					},
				}

				exists, err := bucketExists(mockClient, "bucketname")
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
								{Name: aws.String("anotherbucket")},
							},
						}, nil
					},
				}
				exists, err := bucketExists(mockClient, "bucketname")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		// validating bucket name

		ginkgo.When("The Bucket name is not between 3 and 63 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketNameCompliance("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name must be between 3 and 63 characters long"))
			})
		})

		ginkgo.When("The Bucket name starts with restricted prefixes", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketNameCompliance("xn--bucket")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name cannot start with restricted prefixes (xn-- or sthree-)"))
			})
		})

		ginkgo.When("The Bucket name ends with restricted suffixes", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketNameCompliance("bucket-s3alias")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name cannot end with restricted suffixes (-s3alias or --ol-s3)"))
			})
		})

		ginkgo.When("The Bucket name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketNameCompliance("bucket!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("bucket name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number"))
			})
		})

		ginkgo.When("The Bucket name contains two adjacent periods", func() {
			ginkgo.It("should return an error", func() {
				isValid, _ := checkBucketNameCompliance("bucket..name")
				gomega.Expect(isValid).To(gomega.BeFalse())
			})
		})

		ginkgo.When("The Bucket name is formatted as an IP address", func() {
			ginkgo.It("should return an error", func() {
				isValid, _ := checkBucketNameCompliance("192.168.1.1")
				gomega.Expect(isValid).To(gomega.BeFalse())
			})
		})
	})
})
