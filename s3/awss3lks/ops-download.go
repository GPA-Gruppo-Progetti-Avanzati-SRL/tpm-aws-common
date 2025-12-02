package awss3lks

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

type BlobRange struct {
	Start int64
	End   int64
}

func (br BlobRange) IsZero() bool {
	return br.Start == 0 && br.End == 0
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (lks *LinkedService) DownloadFile(ctx context.Context, bucketName string, objectKey string, blobRange BlobRange) ([]byte, error) {
	const semLogContext = "aws-s3-lks::download-file"

	var rng string
	if !blobRange.IsZero() {
		rng = fmt.Sprintf("bytes=%d-%d", blobRange.Start, blobRange.End)
	}

	result, err := lks.Cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Range:  aws.String(rng),
	})

	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Error().Err(err).Str("bucket", bucketName).Str("object-key", objectKey).Msg(semLogContext + " - not found")
			err = noKey
		} else {
			log.Error().Err(err).Str("bucket", bucketName).Str("object-key", objectKey).Msg(semLogContext)
		}
		return nil, err
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Error().Err(err).Str("bucket", bucketName).Str("object-key", objectKey).Msg(semLogContext)
	}

	return body, err
}
