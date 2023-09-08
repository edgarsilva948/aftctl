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
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// MockS3Client is a mock implementation of an S3 client for testing.
type MockS3Client struct {
	s3iface.S3API
	ListBucketsFunc           func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	CreateBucketFunc          func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	WaitUntilBucketExistsFunc func(*s3.HeadBucketInput) error
	PutPublicAccessBlockFunc  func(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error)
	PutBucketPolicyFunc       func(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)
	PutBucketTaggingFunc      func(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)
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

// PutPublicAccessBlock is a mock implementation of the PutPublicAccessBlock method.
func (m *MockS3Client) PutPublicAccessBlock(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error) {
	if m.PutPublicAccessBlockFunc != nil {
		return m.PutPublicAccessBlockFunc(input)
	}
	return &s3.PutPublicAccessBlockOutput{}, nil
}

// PutBucketPolicy is a mock implementation of the PutBucketPolicy method.
func (m *MockS3Client) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	return m.PutBucketPolicyFunc(input) // Use the custom function field
}

// PutBucketTagging is a mock implementation of the PutBucketTagging method.
func (m *MockS3Client) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	return m.PutBucketTaggingFunc(input) // Use the custom function field
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
				}
				ensure, err := EnsureS3BucketExists(mockClient, "another-bucket", "000000000000", "test-kms-key-id", "codeBuildRole")
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
					PutBucketTaggingFunc: func(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
						return &s3.PutBucketTaggingOutput{}, nil
					},
					PutBucketPolicyFunc: func(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
						return nil, nil
					},
				}
				ensure, err := EnsureS3BucketExists(mockClient, "new-bucket", "000000000000", "test-kms-key-id", "codeBuildRole")
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
				ensure, err := EnsureS3BucketExists(mockClient, "failed-bucket", "000000000000", "test-kms-key-id", "codeBuildRole")
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
				ensure, err := EnsureS3BucketExists(mockClient, "existing-bucket", "000000000000", "test-kms-key-id", "codeBuildRole")
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("AWS WaitUntilBucketExists error"))
			})
		})

		ginkgo.When("PutPublicAccessBlock call fails", func() {
			ginkgo.It("should return an error", func() {
				mockClient := &MockS3Client{
					CreateBucketFunc: func(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
						return &s3.CreateBucketOutput{}, nil
					},
					ListBucketsFunc: func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
						return &s3.ListBucketsOutput{
							Buckets: []*s3.Bucket{},
						}, nil
					},
					WaitUntilBucketExistsFunc: func(input *s3.HeadBucketInput) error {
						return nil
					},
					PutPublicAccessBlockFunc: func(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error) {
						return nil, errors.New("PutPublicAccessBlock failed") // return failure for this method
					},
					PutBucketPolicyFunc: func(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
						return nil, nil
					},
				}

				success, err := EnsureS3BucketExists(mockClient, "validBucketName", "validAftManagementAccountId", "validKmsKeyID", "codeBuildRole")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).To(gomega.BeFalse())
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

	ginkgo.Context("testing the checkIfS3ClientIsProvided", func() {
		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				ensure, err := checkIfS3ClientIsProvided(nil)
				gomega.Expect(ensure).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("S3Client is not provided"))
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
