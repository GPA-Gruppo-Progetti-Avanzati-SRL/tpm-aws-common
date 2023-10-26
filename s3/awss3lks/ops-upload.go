package awss3lks

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

func (lks *LinkedService) UploadBuffer(bucketName string, objectPath, name string, data []byte, contentType string, public bool) (BlobInfo, error) {
	const semLogContext = "aws-s3-lks::upload-file"
	reader := bytes.NewReader(data)

	blob, err := lks.upload(lks.Cli, bucketName, objectPath, name, reader, contentType, public)
	if err != nil {
		log.Error().Err(err).Bool("public", public).Str("bucket-name", bucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
	}

	return blob, err
}

func (lks *LinkedService) UploadFile(bucketName string, objectPath, name string, fileName string, contentType string, public bool) (BlobInfo, error) {
	const semLogContext = "aws-s3-lks::upload-file"
	var blob BlobInfo
	file, err := os.Open(fileName)
	if err != nil {
		log.Error().Err(err).Bool("public", public).Str("file-name", fileName).Str("bucket-name", bucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
	} else {
		defer file.Close()
		blob, err = lks.upload(lks.Cli, bucketName, objectPath, name, file, contentType, public)
		if err != nil {
			log.Error().Err(err).Bool("public", public).Str("file-name", fileName).Str("bucket-name", bucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
		}
	}
	return blob, err
}

func (lks *LinkedService) upload(cli *s3.Client, bucketName string, objectPath, name string, reader io.Reader, contentType string, public bool) (BlobInfo, error) {

	const semLogContext = "aws-s3-lks::upload"

	var err error
	blobName, err := wellFormBlobName(name)
	if err != nil {
		log.Error().Err(err).Bool("public", public).Str("bucket-name", bucketName).Str("name", name).Str("object-path", objectPath).Msg(semLogContext)
		return BlobInfo{}, err
	}

	objectKey := blobName
	if objectPath != "" {
		objectKey = filepath.Join(objectPath, blobName)
	}

	blob := BlobInfo{
		Name:      objectKey,
		Container: bucketName,
		Public:    public,
	}

	acl := types.ObjectCannedACL("")
	if public {
		acl = types.ObjectCannedACLPublicRead
	}

	respOut, err := cli.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(blob.Name),
		Body:        reader,
		ACL:         acl,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Error().Err(err).Bool("public", public).Str("bucket-name", bucketName).Str("object-key", objectKey).Msg(semLogContext)
		return blob, err
	}

	if respOut.VersionId != nil {
		blob.Version = *respOut.VersionId
	}

	return blob, nil
}
