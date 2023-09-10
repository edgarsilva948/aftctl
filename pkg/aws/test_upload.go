package aws

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// PutObject is a mock implementation of the PutObject method.
func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunc(input)
}

var _ = ginkgo.Describe("UploadToS3", func() {
	var mockS3Client *MockS3Client
	var bucketName, bucketKey, fileName string

	ginkgo.BeforeEach(func() {
		mockS3Client = &MockS3Client{
			PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		}
		bucketName = "test-bucket"
		bucketKey = "test-bucket-key"
		fileName = "test-filename"
		ReadFile = os.ReadFile
	})

	ginkgo.Context("when the file is successfully read and uploaded", func() {
		ginkgo.It("should not return an error", func() {
			ReadFile = func(filename string) ([]byte, error) {
				return []byte("fake file content"), nil
			}
			err := UploadToS3(mockS3Client, bucketName, bucketKey, fileName)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})

	ginkgo.When("there is an error reading the file", func() {
		ginkgo.It("should return an error", func() {
			fileName = "nonexistent-file"

			err := UploadToS3(mockS3Client, bucketName, bucketKey, fileName)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.When("there is an error uploading the file", func() {
		ginkgo.It("should return an error", func() {
			mockS3Client.PutObjectFunc = func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				return nil, errors.New("upload error")
			}

			err := UploadToS3(mockS3Client, bucketName, bucketKey, fileName)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
