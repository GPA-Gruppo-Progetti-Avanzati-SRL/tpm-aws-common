package awss3lks

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

type LinkedService struct {
	cfg Config
	Cli *s3.Client
}

func NewLinkedServiceWithConfig(cfg Config) (*LinkedService, error) {
	credProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")

	var serviceClient *s3.Client
	if cfg.UseSharedAWSConfig {

		transport := &http.Transport{
			MaxIdleConns:        100,              // Numero totale di connessioni idle consentite
			MaxIdleConnsPerHost: 100,              // Numero massimo di connessioni per host
			IdleConnTimeout:     90 * time.Second, // Tempo di timeout per connessioni idle
		}

		// Crea un client HTTP personalizzato con il trasporto configurato
		httpClient := &http.Client{
			Transport: transport,
		}

		s3Cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(httpClient))
		if err != nil {
			return nil, err
		}

		serviceClient = s3.NewFromConfig(s3Cfg)
	} else {
		baseEndpoint := cfg.Endpoint
		serviceClient = s3.New(s3.Options{
			BaseEndpoint: &baseEndpoint,
			Credentials:  credProvider,
			Region:       cfg.Region,
		})
	}

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

func (lks *LinkedService) ContainerPublicUrl(bucketName string) string {
	if lks.cfg.PublicEndpoint != "" {
		url := strings.Replace(lks.cfg.PublicEndpoint, "{cnt}", bucketName, -1)
		url = strings.Replace(url, "{region}", lks.cfg.Region, -1)
		return url
	}

	return lks.cfg.PublicEndpoint
}

func (lks *LinkedService) BlobPublicUrl(bucketName string, blobName string) string {
	ep := strings.TrimSuffix(lks.ContainerPublicUrl(bucketName), "/")
	var sb strings.Builder
	sb.WriteString(ep)
	sb.WriteString("/")
	sb.WriteString(blobName)
	return sb.String()
}

func (lks *LinkedService) GetBucketConfig4Map(category string, m map[string]interface{}) (string, string, error) {

	const semLogContext = "aws-s3-linked-service::get-bucket-cfg-4-map"

	var err error
	if len(lks.cfg.BucketConfig) == 0 {
		err = errors.New("no bucket config defined")
		log.Error().Err(err).Msg(semLogContext)
		return "", "", err
	}

	for _, o := range lks.cfg.BucketConfig {
		if o.Category == "*" || o.Category == category {
			cnt, path := resolveBucketTemplates(o.Bucket, o.Path, m)
			return cnt, path, nil
		}
	}

	err = fmt.Errorf("no bucket config found for %s", category)
	log.Error().Err(err).Msg(semLogContext)
	return "", "", err
}

func resolveBucketTemplates(bucket string, path string, m map[string]interface{}) (string, string) {

	if strings.Contains(bucket, "{") || strings.Contains(path, "{") {
		for n, v := range m {
			s := fmt.Sprintf("{%s}", n)
			bucket = strings.Replace(bucket, s, fmt.Sprint(v), -1)
			path = strings.Replace(path, s, fmt.Sprint(v), -1)
		}
	}

	return bucket, path
}
