package logs

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func ShowBuckets(ctx context.Context, client *minio.Client) (bool, error){
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return false, err
	}

	fmt.Println("Buckets: ")
	for _, bucket := range buckets {
		fmt.Println("-----------")
		fmt.Println(bucket.Name)
	}
	fmt.Println("-----------")

	return true, nil
}