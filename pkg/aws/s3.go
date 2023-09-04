/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketExists checks if a given S3 bucket exists.
func BucketExists(client S3Client, bucketName string) (bool, error) {
	// Check if the client is nil
	if client == nil {
		return false, fmt.Errorf("S3Client is not provided")
	}

	input := &s3.ListBucketsInput{}

	output, err := client.ListBuckets(input)
	if err != nil {
		return false, fmt.Errorf("failed to list S3 buckets: %w", err)
	}

	for _, bucket := range output.Buckets {
		if *bucket.Name == bucketName {
			return true, nil
		}
	}

	return false, nil
}

// CreateS3Bucket creates a new S3 bucket with the given name.
func CreateS3Bucket(bucketName string, kmsKeyID string, region string) error {
	// new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	// Ativar o versionamento
	// _, err = svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
	// 	Bucket: aws.String(bucketName),
	// 	VersioningConfiguration: &s3.VersioningConfiguration{
	// 		Status: aws.String("Enabled"),
	// 	},
	// })

	if err != nil {
		return err
	}

	// Bloqueio de acesso público
	_, err = svc.PutPublicAccessBlock(&s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls:   aws.Bool(true),
			BlockPublicPolicy: aws.Bool(true),
		},
	})

	if err != nil {
		return err
	}

	// Configuração do ciclo de vida
	// _, err = svc.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
	// 	Bucket: aws.String(bucketName),
	// 	LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
	// 		Rules: []*s3.LifecycleRule{
	// 			{
	// 				Status: aws.String("Enabled"),
	// 				// Adicione suas regras aqui
	// 			},
	// 		},
	// 	},
	// })

	if err != nil {
		return err
	}

	// Criptografia KMS
	_, err = svc.PutBucketEncryption(&s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucketName),
		ServerSideEncryptionConfiguration: &s3.ServerSideEncryptionConfiguration{
			Rules: []*s3.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &s3.ServerSideEncryptionByDefault{
						KMSMasterKeyID: aws.String(kmsKeyID),
						SSEAlgorithm:   aws.String("aws:kms"),
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
