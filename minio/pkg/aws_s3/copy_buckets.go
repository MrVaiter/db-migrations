package aws_s3

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CopyPredicate func(string) bool
type ClearPredicate func(string) bool

func (from *Client) CopyBucketsWithFilter(ctx context.Context,to *Client, filter CopyPredicate) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})


	log.Print("Starting copying buckets from ", from.EndpointURL().Host, " to ", to.EndpointURL().Host)

	buckets, err := from.ListBuckets(ctx)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, bucket := range buckets {
		if filter(bucket.Name) {
			buffer.WriteString(bucket.Name)
			
			exist, err := to.BucketExists(context.Background() ,bucket.Name)
			if  !exist {
				err = to.MakeBucket(ctx, bucket.Name, minio.MakeBucketOptions{})
				if err != nil {
					buffer.WriteString(" \tfailed")
					log.Error().Msg(buffer.String())
					buffer.Reset()
					return err
				}

				buffer.WriteString(" \tsuccess")
			} else {
				buffer.WriteString(" \talready exists")
			}
		}

		log.Info().Msg(buffer.String())
		buffer.Reset()
	}

	return nil
}

func (from *Client) CopyBucketsWithSuffix(ctx context.Context,to *Client, suffix string) error {
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
