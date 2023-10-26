package awss3lks

import (
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
)

type LinkedService struct {
	cfg Config
	Cli *s3.Client
}

func NewLinkedServiceWithConfig(cfg Config) (*LinkedService, error) {
	credProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")

	baseEndpoint := cfg.Endpoint
	serviceClient := s3.New(s3.Options{
		BaseEndpoint: &baseEndpoint,
		Credentials:  credProvider,
		Region:       cfg.Region,
	})

	lks := &LinkedService{cfg: cfg, Cli: serviceClient}
	return lks, nil
}

func NewLinkedService(name string, opts ...Option) (*LinkedService, error) {
	cfg := Config{Name: name}

	for _, o := range opts {
		o(&cfg)
	}

	return NewLinkedServiceWithConfig(cfg)
}

func (lks *LinkedService) PublicUrlOf(bucketName string) string {
	if lks.cfg.PublicEndpoint != "" {
		url := strings.Replace(lks.cfg.PublicEndpoint, "{cnt}", bucketName, -1)
		url = strings.Replace(url, "{region}", lks.cfg.Region, -1)
		return url
	}

	return lks.cfg.PublicEndpoint
}
