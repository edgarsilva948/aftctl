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
	ListBucketsFunc           func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	CreateBucketFunc          func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	WaitUntilBucketExistsFunc func(*s3.HeadBucketInput) error
}

// ListBuckets is a mock implementation of the ListBuckets method.
func (m *MockS3Client) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(input)
}

// CreateBucket is a mock implementation of the CreateBucket method.
func (m *MockS3Client) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketFunc(input)
}

// WaitUntilBucketExists is a mock implementation of the WaitUntilBucketExists method.
func (m *MockS3Client) WaitUntilBucketExists(input *s3.HeadBucketInput) error {
	return m.WaitUntilBucketExistsFunc(input)
}

var _ = ginkgo.Describe("Interacting with the S3 API", func() {

	ginkgo.Context("testing the EnsureS3bucketExists function", func() {

		ginkgo.When("bucket already exists", func() {
			ginkgo.It("should return an success", func() {

				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("another-bucket")},
							},
						}, nil
					},
					CreateBucketFunc: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
						return &s3.CreateBucketOutput{}, nil
					},
					WaitUntilBucketExistsFunc: func(input *s3.HeadBucketInput) error {
						return nil
					},
				}
				ensure, err := EnsureS3BucketExists(mockClient, "another-bucket", "test-kms-key-id")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("bucket doesn't exists", func() {
			ginkgo.It("should create the bucket", func() {

				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("")},
							},
						}, nil
					},
					CreateBucketFunc: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
						return &s3.CreateBucketOutput{}, nil
					},
					WaitUntilBucketExistsFunc: func(input *s3.HeadBucketInput) error {
						return nil
					},
				}
				ensure, err := EnsureS3BucketExists(mockClient, "new-bucket", "test-kms-key-id")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("bucket creation fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					CreateBucketFunc: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
						return nil, errors.New("AWS create bucket error")
					},
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{},
						}, nil
					},
					WaitUntilBucketExistsFunc: func(input *s3.HeadBucketInput) error {
						return nil
					},
				}
				ensure, err := EnsureS3BucketExists(mockClient, "failed-bucket", "test-kms-key-id")
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("AWS create bucket error"))
			})
		})

		ginkgo.When("Wait Until Bucket Exists fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					CreateBucketFunc: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
						return nil, errors.New("AWS WaitUntilBucketExists error")
					},
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{},
						}, nil
					},
					WaitUntilBucketExistsFunc: func(input *s3.HeadBucketInput) error {
						return nil
					},
				}
				ensure, err := EnsureS3BucketExists(mockClient, "existing-bucket", "test-kms-key-id")
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("AWS WaitUntilBucketExists error"))
			})
		})

	})

	ginkgo.Context("testing the bucketExists function", func() {

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
								{Name: aws.String("another-bucket")},
							},
						}, nil
					},
				}
				exists, err := bucketExists(mockClient, "bucketname")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("AWS returns an empty bucket list", func() {
			ginkgo.It("should return false", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{}, nil
					},
				}
				exists, err := bucketExists(mockClient, "non-existent-bucket")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("testing the checkIfClientIsProvided", func() {
		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("S3Client is not provided"))
			})
		})
	})

	ginkgo.Context("testing the bucketExists", func() {
		ginkgo.When("if bucket name is not provided", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("another-bucket")},
							},
						}, nil
					},
				}
				ensure, err := bucketExists(mockClient, "")
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Context("testing the checkIfBucketNameIsProvided function", func() {
		ginkgo.When("bucketName is not provided", func() {
			ginkgo.It("should return an error", func() {

				check, err := checkIfBucketNameIsProvided("")
				gomega.Expect(check).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name is not provided"))
			})
		})
	})

	ginkgo.Context("testing the checkBucketNameCompliance", func() {
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

		ginkgo.When("The Bucket name is valid", func() {
			ginkgo.It("should return true", func() {
				isValid, err := checkBucketNameCompliance("valid-bucket-name")
				gomega.Expect(isValid).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("testing the checkIfBucketExists function", func() {
		ginkgo.When("if the bucket already exists", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{
								{Name: aws.String("one-bucket")},
							},
						}, nil
					},
				}
				ensure, err := checkIfBucketExists(mockClient, "one-bucket")
				gomega.Expect(ensure).To(gomega.BeTrue())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.When("list operation fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return nil, errors.New("AWS S3 list operation error")
					},
				}
				exists, err := checkIfBucketExists(mockClient, "any-bucket")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("failed to list S3 buckets: AWS S3 list operation error"))
			})
		})
	})
})
