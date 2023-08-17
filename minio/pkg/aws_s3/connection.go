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

	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Client:  minioClient,
		buckets: make([]string, 0, len(buckets)),
	}

	for _, bucket := range buckets {
		client.buckets = append(client.buckets, bucket.Name)
	}

	return client, nil
}
