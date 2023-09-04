/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aft

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

// MockS3Client is a mock implementation of an S3 client.
// It provides a way to specify the behavior of the ListBuckets method
// via the ListBucketsFunc field.
type MockS3Client struct {
	ListBucketsFunc func(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

// var _ = ginkgo.Describe("Interacting with the S3 API", func() {

// 	ginkgo.Context("Validating the CreateStateBucketInAftAccount() with the AWS sdk", func() {

// 		ginkgo.When("S3Client is not provided", func() {
// 			ginkgo.It("should return an error", func() {
// 				exists, err := CreateStateBucketInAftAccount(nil, "bucket-name")
// 				gomega.Expect(exists).To(gomega.BeFalse())
// 				gomega.Expect(err).To(gomega.MatchError("error checking if bucket exists: S3Client is not provided"))
// 			})
// 		})

// 	})
// })

// var _ = ginkgo.Describe("Validating the bucket name provided by the user", func() {

// 	ginkgo.Context("Testing the provided name using the func checkBucketName()", func() {

// })
