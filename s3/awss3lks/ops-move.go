package awss3lks

import (
	"context"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func (lks *LinkedService) Copy(fromBucketName, toBucketName string, objectPath, name string) error {
	const semLogContext = "aws-s3-lks::delete-blob"

	err := lks.copyMove(lks.Cli, fromBucketName, toBucketName, objectPath, name, false)
	if err != nil {
		log.Error().Err(err).Str("bucket-name", fromBucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
	}

	return err
}

func (lks *LinkedService) Move(fromBucketName, toBucketName string, objectPath, name string) error {
	const semLogContext = "aws-s3-lks::delete-blob"

	err := lks.copyMove(lks.Cli, fromBucketName, toBucketName, objectPath, name, true)
	if err != nil {
		log.Error().Err(err).Str("bucket-name", fromBucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
	}

	return err
}

func (lks *LinkedService) copyMove(cli *s3.Client, fromBucketName, toBucketName string, objectPath, name string, isMove bool) error {
	const semLogContext = "aws-s3-lks::delete"

	objectKey := path.Join(objectPath, name)
	resp, err := cli.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(toBucketName),
		Key:        aws.String(objectKey),
		CopySource: aws.String(url.QueryEscape(path.Join(fromBucketName, objectKey))),
	})

	if err != nil {
		log.Error().Err(err).Str("from-bucket-name", fromBucketName).Str("to-bucket-name", toBucketName).Str("object-key", objectKey).Msg(semLogContext)
		return err
	} else {
		log.Info().Str("from-bucket-name", fromBucketName).Str("to-bucket-name", toBucketName).Str("object-key", objectKey).Interface("resp", resp).Msg(semLogContext)
	}

	if isMove {
		return lks.delete(cli, fromBucketName, objectPath, name)
	}

	return nil
}
