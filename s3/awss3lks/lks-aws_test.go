package awss3lks_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-aws-common/s3/awss3lks"
	"github.com/stretchr/testify/require"
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

	//AWS_TargetContainer = "delta-tests"
	//AWS_Region          = "eu-north-1"
	//AWS_PublicEndpoint  = "https://{cnt}.s3.{region}.amazonaws.com/"
	//MAXObjectsToUpload  = 1000
	//DoRetrieve          = false

	AWS_TargetContainer = "gpagroup-dev-opem-flussi"
	AWS_Region          = "eu-central-1"
	AWS_PublicEndpoint  = "https://{cnt}.s3.{region}.amazonaws.com/"
	MAXObjectsToUpload  = 1000
	DoRetrieve          = false
)

func TestClientAWS(t *testing.T) {

	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         os.Getenv(AWSS3RegionEnvVarName),
		PublicEndpoint: "", // AWS_PublicEndpoint,
	}

	// cfg.Endpoint = "https://s3.eu-central-1.amazonaws.com"
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

func TestDownload(t *testing.T) {
	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         os.Getenv(AWSS3RegionEnvVarName),
		PublicEndpoint: "", // AWS_PublicEndpoint,
	}

	require.True(t, cfg.Endpoint != "")
	require.True(t, cfg.AccessKey != "")
	require.True(t, cfg.SecretKey != "")

	lks, err := awss3lks.NewLinkedServiceWithConfig(cfg)
	require.NoError(t, err)

	b, err := lks.DownloadFile(context.Background(), "opem-warehousereportrange-10000", "range/cms_10000_PLD001_L6_AM_RUOWR001_1763486778", awss3lks.BlobRange{End: 10})
	require.NoError(t, err)

	t.Log(string(b))
}

func TestMove(t *testing.T) {
	cfg := awss3lks.Config{
		Endpoint:       os.Getenv(AWSCommonEndpointEnvVarName),
		AccessKey:      os.Getenv(AWSCommonAccessKeyEnvVarName),
		SecretKey:      os.Getenv(AWSCommonSecretKeyEnvVarName),
		Region:         os.Getenv(AWSS3RegionEnvVarName),
		PublicEndpoint: "", // AWS_PublicEndpoint,
	}

	require.True(t, cfg.Endpoint != "")
	require.True(t, cfg.AccessKey != "")
	require.True(t, cfg.SecretKey != "")

	lks, err := awss3lks.NewLinkedServiceWithConfig(cfg)
	require.NoError(t, err)

	err = lks.Move("opem-from-istituto-10000-queue", "opem-from-istituto-10000", "", "card_assegnate_20251203.xml")
	require.NoError(t, err)
}
