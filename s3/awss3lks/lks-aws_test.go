package awss3lks_test

import (
	_ "embed"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	// //TargetContainer     = "r3ds9-s3-test-cnt"
	// TargetContainer     = "r3ds9-s3-user-cnt"
	// Region              = "gra"
	// PublicEndpoint      = "https://{cnt}.s3.{region}.io.cloud.ovh.net/"
	// DropContainerOnExit = false
	//
	// AWSCommonAccessKeyEnvVarName = "AWSCOMMON_ACCESSKEY"
	// AWSCommonSecretKeyEnvVarName = "AWSCOMMON_SECRETKEY"
	// AWSCommonEndpointEnvVarName  = "AWSCOMMON_ENDPOINT"

	AWS_TargetContainer = "delta-tests"
	AWS_Region          = "eu-north-1"
	AWS_PublicEndpoint  = "https://{cnt}.s3.{region}.amazonaws.com/"
	MAXObjectsToUpload  = 1000
	DoRetrieve          = false
)

func TestClientAWS(t *testing.T) {

	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         AWS_Region,
		PublicEndpoint: "", // AWS_PublicEndpoint,
	}

	cfg.Endpoint = "https://s3.eu-north-1.amazonaws.com"
	require.True(t, cfg.Endpoint != "")
	require.True(t, cfg.AccessKey != "")
	require.True(t, cfg.SecretKey != "")

	lks, _ := awss3lks.NewLinkedServiceWithConfig(cfg)
	buckets, err := lks.ListBuckets()
	require.NoError(t, err)
	for _, b := range buckets {
		t.Log(b)
	}

	blob, err := lks.UploadBuffer(AWS_TargetContainer, "", "user-photo.jpg", img, "image/jpeg", false)
	require.NoError(t, err)
	t.Log(blob)

	blob, err = lks.UploadFile(AWS_TargetContainer, "p1", "zucca  mia bella zucca!.txt", "config.go", "text/plain", false)
	require.NoError(t, err)
	t.Log(blob)

	for i := 0; i < MAXObjectsToUpload; i++ {
		blob, err = lks.UploadFile(AWS_TargetContainer, "p1/p2", fmt.Sprintf("zucca  mia bella zucca!-%02d.txt", i), "config.go", "text/plain", false)
		require.NoError(t, err)
		t.Log(blob)
	}

	if DoRetrieve {
		pr := lks.NewListObjectsPager(AWS_TargetContainer, 5)
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
}
