/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"github.com/aws/aws-sdk-go/service/s3"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// MockS3Client is a mock implementation of an S3 client.
// It provides a way to specify the behavior of the ListBuckets method
// via the ListBucketsFunc field.
type MockS3Client struct {
	ListBucketsFunc func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

var _ = ginkgo.Describe("Interacting with the S3 API", func() {

	ginkgo.Context("Validating the CreateStateBucketInAftAccount() with the AWS sdk", func() {

		ginkgo.When("S3Client is not provided", func() {
			ginkgo.It("should return an error", func() {
				exists, err := CreateStateBucketInAftAccount(nil, "bucket-name")
				gomega.Expect(exists).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("error checking if bucket exists: S3Client is not provided"))
			})
		})

	})
})

var _ = ginkgo.Describe("Validating the bucket name provided by the user", func() {

	ginkgo.Context("Testing the provided name using the func checkBucketName()", func() {

		ginkgo.When("The Bucket name is not between 3 and 63 characters long", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketName("ab")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name must be between 3 and 63 characters long"))
			})
		})

		ginkgo.When("The Bucket name starts with restricted prefixes", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketName("xn--bucket")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name cannot start with restricted prefixes (xn-- or sthree-)"))
			})
		})

		ginkgo.When("The Bucket name ends with restricted suffixes", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketName("bucket-s3alias")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err).To(gomega.MatchError("bucket name cannot end with restricted suffixes (-s3alias or --ol-s3)"))
			})
		})

		ginkgo.When("The Bucket name contains invalid characters", func() {
			ginkgo.It("should return an error", func() {
				isValid, err := checkBucketName("bucket!")
				gomega.Expect(isValid).To(gomega.BeFalse())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("bucket name can only consist of lowercase letters, numbers, and hyphens, and must begin and end with a letter or number"))
			})
		})

		ginkgo.When("The Bucket name contains two adjacent periods", func() {
			ginkgo.It("should return an error", func() {
				isValid, _ := checkBucketName("bucket..name")
				gomega.Expect(isValid).To(gomega.BeFalse())
			})
		})

		ginkgo.When("The Bucket name is formatted as an IP address", func() {
			ginkgo.It("should return an error", func() {
				isValid, _ := checkBucketName("192.168.1.1")
				gomega.Expect(isValid).To(gomega.BeFalse())
			})
		})
	})
})
