package aws_s3

import (
	"context"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CopyPredicate func(string) bool
type ClearPredicate func(string) bool

func (from *Client) CopyBucketsWithFilter(ctx context.Context, to *Client, filter CopyPredicate) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	buckets, err := from.ListBuckets(ctx)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if filter(bucket.Name) {

			exist, err := to.BucketExists(context.Background(), bucket.Name)
			if err != nil {
				return err
			}
			
			if !exist {
				err = to.MakeBucket(ctx, bucket.Name, minio.MakeBucketOptions{})
				if err != nil {
					return err
				}

				log.Debug().Str("bucket", bucket.Name).Msg("Success")
			} else {
				log.Debug().Str("bucket", bucket.Name).Msg("Already exists")
			}
		}
	}

	return nil
}

func (from *Client) CopyBucketsWithSuffix(ctx context.Context, to *Client, suffix string) error {
	return from.CopyBucketsWithFilter(ctx, to, func(name string) bool {
		return strings.Contains(name, suffix)
	})
}

func (client *Client) ClearWithFilter(ctx context.Context, filter ClearPredicate) error {

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if filter(bucket.Name) {
			err = client.RemoveBucketWithOptions(ctx, bucket.Name, minio.RemoveBucketOptions{ForceDelete: true})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *Client) ClearWithSuffix(ctx context.Context, suffix string) error {
	return client.ClearWithFilter(ctx, func(name string) bool {
		return strings.Contains(name, suffix)
	})
}
