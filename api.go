package titans3

import (
	"github.com/Filecoin-Titan/titan-sdk-go-s3/config"
	"github.com/Filecoin-Titan/titan-sdk-go-s3/titan"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"log"
)

type API interface {
	GetLocalMinioAPI() (*minio.Client, error)
}

type Client struct {
	cfg   config.Config
	titan *titan.Service
}

func New(opts ...config.Option) *Client {
	cfg := config.DefaultOption()

	for _, opt := range opts {
		opt(&cfg)
	}

	s, err := titan.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{
		cfg:   cfg,
		titan: s,
	}
}

func (c *Client) GetLocalMinioAPI() (*minio.Client, error) {
	schedulerURL, err := c.titan.GetScheduler()
	if err != nil {
		return nil, err
	}

	local, err := c.titan.GetLocalMinioEndpoint(schedulerURL)
	if err != nil {
		return nil, errors.Errorf("get local minio endpoint failed: %v", err)
	}

	// Initialize minio client object.
	minioClient, err := minio.New(local.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(local.AccessKeyID, local.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, errors.Errorf("create minio client failed: %v", err)
	}

	return minioClient, nil
}

var _ API = (*Client)(nil)
