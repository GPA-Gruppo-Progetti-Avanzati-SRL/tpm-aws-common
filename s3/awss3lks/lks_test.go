package awss3lks_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
)

const (
	//TargetContainer     = "r3ds9-s3-test-cnt"
	TargetContainer     = "r3ds9-s3-user-cnt"
	Region              = "gra"
	PublicEndpoint      = "https://{cnt}.s3.{region}.io.cloud.ovh.net/"
	DropContainerOnExit = false

	AWSCommonAccessKeyEnvVarName = "AWS_S3_ACCESS_KEY"
	AWSCommonSecretKeyEnvVarName = "AWS_S3_SECRET_KEY"
	AWSCommonEndpointEnvVarName  = "AWS_S3_ENDPOINT"
	AWSS3RegionEnvVarName        = "AWS_S3_REGION"
)

//go:embed test-image.jpg
var img []byte

func TestClient(t *testing.T) {

	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         Region,
		PublicEndpoint: PublicEndpoint,
	}

	require.True(t, cfg.Endpoint != "")
	require.True(t, cfg.AccessKey != "")
	require.True(t, cfg.SecretKey != "")

	lks, _ := awss3lks.NewLinkedServiceWithConfig(cfg)
	buckets, err := lks.ListBuckets()
	require.NoError(t, err)
	for _, b := range buckets {
		t.Log(b)
	}

	blob, err := lks.UploadBuffer(TargetContainer, "", "user-photo.jpg", img, "image/jpeg", true)
	require.NoError(t, err)
	t.Log(blob)

	blob, err = lks.UploadFile(TargetContainer, "p1", "zucca  mia bella zucca!.txt", "config.go", "text/plain", true)
	require.NoError(t, err)
	t.Log(blob)

	for i := 0; i < 20; i++ {
		blob, err = lks.UploadFile(TargetContainer, "p1/p2", fmt.Sprintf("zucca  mia bella zucca!-%02d.txt", i), "config.go", "text/plain", true)
		require.NoError(t, err)
		t.Log(blob)
	}

	pr := lks.NewListObjectsPager(TargetContainer, 5)
	objs, eol, err := pr.Next()
	require.NoError(t, err)
	t.Log("Page:", pr.Page(), "Num objs:", len(objs))
	for _, o := range objs {
		logBlobInfo(t, o)
	}
	for eol {
		objs, eol, err = pr.Next()
		require.NoError(t, err)
		t.Log("Page:", pr.Page(), "Num objs:", len(objs))
		for _, o := range objs {
			logBlobInfo(t, o)
		}
	}
}

func logBlobInfo(t *testing.T, bl awss3lks.BlobInfo) {
	b, err := json.Marshal(bl)
	require.NoError(t, err)
	t.Log(string(b))

	ep, err := bl.PublicEndpoint()
	require.NoError(t, err)
	t.Log("url: ", ep)
}

func TestWellFormed(t *testing.T) {

	n := "a b c"
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._"

	v := len(alphabet) % 33
	t.Log(v, len(alphabet))

	var out strings.Builder
	var err error
	lenAlphabet := int32(len(alphabet))
	for _, c := range n {
		if strings.Index(alphabet, string(c)) >= 0 {
			_, err = out.WriteString(string(c))
		} else {
			if c > lenAlphabet {
				c = c % lenAlphabet
			} else {
				c = lenAlphabet % c
			}
			_, err = out.WriteString(string(c))
		}

		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log(out.String())
}

func createBucket(cli *s3.Client, name string, region string) error {
	_, err := cli.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		ACL:    types.BucketCannedACLPublicRead,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	if err != nil {
		log.Printf("Couldn't create bucket %v in Region %v. Here's why: %v\n",
			name, region, err)
	}
	return err
}

func uploadFile(cli *s3.Client, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = cli.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(objectKey),
			Body:        file,
			ACL:         types.ObjectCannedACLPublicRead,
			ContentType: aws.String("text/plain"),
			//GrantRead: aws.String("public-read"), // http://acs.amazonaws.com/groups/global/AllUsers,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}
