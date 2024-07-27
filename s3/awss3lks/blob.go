package awss3lks

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type BlobInfo struct {
	publicEndpoint string
	Key            string `mapstructure:"key,omitempty" yaml:"key,omitempty" json:"key,omitempty"`
	Container      string `mapstructure:"bucket,omitempty" yaml:"bucket,omitempty" json:"bucket,omitempty"`
	Version        string `mapstructure:"version,omitempty" yaml:"version,omitempty" json:"version,omitempty"`
	Public         bool   `mapstructure:"public,omitempty" yaml:"public,omitempty" json:"public,omitempty"`
	ETag           string `mapstructure:"etag,omitempty" yaml:"etag,omitempty" json:"etag,omitempty"`
	LastModified   string `mapstructure:"last-modified,omitempty" yaml:"last-modified,omitempty" json:"last-modified,omitempty"`
	Size           int    `mapstructure:"size,omitempty" yaml:"size,omitempty" json:"size,omitempty"`
}

func (bi *BlobInfo) PublicEndpoint() (string, error) {
	if bi.publicEndpoint == "" {
		return "", fmt.Errorf("blob %s has no public endpoint set", bi.Key)
	}

	ep := strings.Replace(strings.TrimSuffix(bi.publicEndpoint, "/"), "{cnt}", bi.Container, 1)
	var sb strings.Builder
	sb.WriteString(ep)
	sb.WriteString("/")
	sb.WriteString(bi.Key)
	return sb.String(), nil
}

func NewBlobInfoFromObject(bucketName string, obj types.Object) BlobInfo {
	const semLogContext = "aws-s3-lks-blob::new-blob-from-object"
	b := BlobInfo{
		Key:          *obj.Key,
		Container:    bucketName,
		Version:      "",
		Public:       false,
		ETag:         *obj.ETag,
		LastModified: obj.LastModified.Format(time.RFC3339Nano),
		Size:         int(*obj.Size), // TODO This can be null as a pointer?
	}

	log.Trace().Interface("owner", obj.Owner).Msg(semLogContext)
	return b
}
