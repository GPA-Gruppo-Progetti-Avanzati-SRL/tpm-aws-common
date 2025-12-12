package awss3lks

import (
	"context"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (lks *LinkedService) Delete(bucketName string, objectPath, name string) error {
	const semLogContext = "aws-s3-lks::delete-blob"

	err := lks.delete(lks.Cli, bucketName, objectPath, name)
	if err != nil {
		log.Error().Err(err).Str("bucket-name", bucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
	}

	return err
}

func (lks *LinkedService) delete(cli *s3.Client, bucketName string, objectPath, name string) error {
	const semLogContext = "aws-s3-lks::delete"

	objectKey := path.Join(objectPath, name)
	resp, err := cli.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Error().Err(err).Str("bucket-name", bucketName).Str("object-key", objectKey).Msg(semLogContext)
		return err
	} else {
		log.Info().Str("bucket-name", bucketName).Str("object-key", objectKey).Interface("resp", resp).Msg(semLogContext)
	}

	return err
}
