package aws_s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

type CopyPredicate func(string) bool
type ClearPredicate func(string) bool

func (from *Client) CopyBucketsWithFilter(to *Client, filter CopyPredicate) error {

	ctx := context.Background()

	fmt.Println("Starting copying buckets from", from.EndpointURL().Host, "to", to.EndpointURL().Host)

	buckets, err := from.ListBuckets(ctx)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if filter(bucket.Name) {
			fmt.Print(bucket.Name)

			if !to.BucketExists(bucket.Name) {
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

func (from *Client) CopyBucketsWithSuffix(to *Client, suffix string) error {
	return from.CopyBucketsWithFilter(to, func(name string) bool {
		return strings.Contains(name, suffix)
	})
}

func (client *Client) BucketExists(bucketName string) bool {

	result, err := client.Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		panic(err)
	}

	return result
}

func (client *Client) ClearWithFilter(filter ClearPredicate) error {

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if filter(bucket.Name) {
			err = client.RemoveBucketWithOptions(context.Background(), bucket.Name, minio.RemoveBucketOptions{ForceDelete: true})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *Client) ClearWithSuffix(suffix string) error {
	return client.ClearWithFilter(func(name string) bool {
		return strings.Contains(name, suffix)
	})
}
