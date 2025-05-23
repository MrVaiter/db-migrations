package aws_s3

import (
	"context"
	"github.com/minio/minio-go/v7"
)

const filesBatchSize = 1000

func (client *Client) ListFiles(ctx context.Context) ([]*FileHandle, error) {
	handles := make([]*FileHandle, 0, filesBatchSize)

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	for _, bucket := range buckets {

		hasMore := true
		continueFrom := ""

		for hasMore {
			objectsCh := client.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
				WithVersions: false,
				WithMetadata: false,
				Recursive:    true,
				MaxKeys:      filesBatchSize,
				StartAfter:   continueFrom,
				UseV1:        true,
			})
			totalAdded := 0

			for object := range objectsCh {
				if object.Err != nil {
					return nil, object.Err
				}

				handles = append(handles, &FileHandle{
					bucketName: bucket.Name,
					fileName:   object.Key,
					size:       object.Size,
				})

				totalAdded += 1
			}

			if totalAdded == 0 {
				hasMore = false
			} else {
				continueFrom = handles[len(handles)-1].fileName
			}
		}
	}

	return handles, nil
}
