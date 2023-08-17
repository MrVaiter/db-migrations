package aws_s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

func CopyBuckets(from *Client, to *Client, filter string) error {

	ctx := context.Background()

	fmt.Println("Starting copying buckets from", from.EndpointURL().Host, "to", to.EndpointURL().Host)

	buckets, err := from.ListBuckets(ctx)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if strings.Contains(bucket.Name, filter) {
			fmt.Print(bucket.Name)

			if !checkBucket(to, bucket.Name) {
				err = to.MakeBucket(ctx, bucket.Name, minio.MakeBucketOptions{})
				if err != nil {
					fmt.Println(" failed")
					return err
				}

				fmt.Println(" success")
			} else {
				fmt.Println(" already exists")
			}
		}
	}

	return nil
}

func checkBucket(client *Client, bucketName string) bool {

	result, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		panic(err)
	}

	return result
}

func clear(client *Client) error {

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if strings.Contains(bucket.Name, "-test-") {
			err = client.RemoveBucketWithOptions(context.Background(), bucket.Name, minio.RemoveBucketOptions{ForceDelete: true})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
