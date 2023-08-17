package aws_s3

import "github.com/minio/minio-go/v7"

type Client struct {
	*minio.Client
}