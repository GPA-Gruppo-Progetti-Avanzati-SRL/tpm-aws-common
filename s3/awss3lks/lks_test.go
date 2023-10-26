package awss3lks_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"tpm-aws-common/s3/awss3lks"
)

const (
	TargetContainer     = "r3ds9-s3-test-cnt"
	DropContainerOnExit = false

	AWSCommonAccessKeyEnvVarName = "AWSCOMMON_ACCESSKEY"
	AWSCommonSecretKeyEnvVarName = "AWSCOMMON_SECRETKEY"
	AWSCommonEndpointEnvVarName  = "AWSCOMMON_ENDPOINT"
)

func TestClient(t *testing.T) {

	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         "gra",
		PublicEndpoint: "https://{cnt}.s3.{region}.io.cloud.ovh.net/",
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

	blob, err := lks.UploadFile(TargetContainer, "p1", "zucca  mia bella zucca!.txt", "config.go", "text/plain", true)
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
		t.Log(filepath.Join(lks.PublicUrlOf(TargetContainer), *o.Key))
	}
	for eol {
		objs, eol, err = pr.Next()
		require.NoError(t, err)
		t.Log("Page:", pr.Page(), "Num objs:", len(objs))
		for _, o := range objs {
			t.Log(filepath.Join(lks.PublicUrlOf(TargetContainer), *o.Key))
		}
	}
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
