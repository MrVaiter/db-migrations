package aws_s3

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(ctx context.Context, endpoint, accessKeyID, secretAccessKey, token string) (*Client, error) {
	withSSLValue, ok := os.LookupEnv("WITH_SSL")
	useSSL := ok && withSSLValue != "false" && withSSLValue != "False"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, token),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	client := &Client{
		Client: minioClient,
	}

	return client, nil
}
