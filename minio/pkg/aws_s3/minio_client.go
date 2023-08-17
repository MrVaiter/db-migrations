package aws_s3

import "github.com/minio/minio-go/v7"

type Client struct {
	*minio.Client
	buckets []string
}

func (c *Client) GetBuckets() []string {
	return c.buckets
}

func (c *Client) GetMinioClient() *minio.Client {
	return c.Client
}
