package awss3lks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

func (lks *LinkedService) ListBuckets() ([]string, error) {
	opts := s3.ListBucketsInput{}
	listOut, err := lks.Cli.ListBuckets(context.TODO(), &opts)
	if err != nil {
		return nil, err
	}

	var res []string
	if listOut != nil {
		for _, b := range listOut.Buckets {
			res = append(res, *b.Name)
		}
	}

	return res, nil
}

type ListObjectsPager struct {
	lks               *LinkedService
	bucketName        string
	maxKeys           int32
	continuationToken string
	pageNumber        int
}

func (pr *ListObjectsPager) Page() int {
	return pr.pageNumber
}

func (lks *LinkedService) NewListObjectsPager(bucketName string, pageSize int) ListObjectsPager {
	return ListObjectsPager{lks: lks, bucketName: bucketName, maxKeys: int32(pageSize)}
}

func (qry *ListObjectsPager) Next() ([]BlobInfo, bool, error) {

	const semLogContext = "aws-s3-lks::list-objects"

	input := s3.ListObjectsV2Input{
		Bucket: aws.String(qry.bucketName),
	}

	if qry.maxKeys > 0 {
		input.MaxKeys = &qry.maxKeys
	}

	if qry.continuationToken != "" {
		input.ContinuationToken = &qry.continuationToken
		qry.pageNumber++
	} else {
		qry.pageNumber = 0
	}

	result, err := qry.lks.Cli.ListObjectsV2(context.TODO(), &input)
	if err != nil {
		log.Error().Err(err).Str("bucket-name", qry.bucketName).Msg(semLogContext)
		return nil, false, err
	}

	var contents []BlobInfo
	contents = qry.adaptToBlobInfo(qry.bucketName, result.Contents)

	if result.NextContinuationToken != nil && result.IsTruncated != nil && *result.IsTruncated {
		qry.continuationToken = *result.NextContinuationToken
	}

	return contents, result.IsTruncated != nil && *result.IsTruncated, err
}

func (qry *ListObjectsPager) adaptToBlobInfo(cnt string, objs []types.Object) []BlobInfo {

	var binfos []BlobInfo
	for _, o := range objs {
		b := NewBlobInfoFromObject(cnt, o)
		b.publicEndpoint = qry.lks.ContainerPublicUrl(cnt)
		binfos = append(binfos, b)
	}
	return binfos
}
