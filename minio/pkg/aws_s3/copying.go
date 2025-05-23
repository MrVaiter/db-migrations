package aws_s3

import (
	"context"
	"io"
	"sync"

	"github.com/minio/minio-go/v7"
)

type copyQueue struct {
	jobs       chan *FileHandle
	fromClient *Client
	toClient   *Client

	ctx    context.Context
	cancel context.CancelFunc
}

func (q *copyQueue) toCopy(jobs []*FileHandle) {
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	go func() {

		for _, job := range jobs {
			q.jobs <- job
			wg.Done()
		}
		
	}()

	wg.Wait()
}

func (q *copyQueue) doCopy(ctx context.Context, results *[]*CopyResult) bool { // TODO: return results ???

	loops := len(q.jobs)

	var currentBucket string
	for i := 0; i < loops; i++ {

		select {
		// if context was canceled.
		case <-q.ctx.Done():
			return true
		// if job received.
		case job := <-q.jobs:

			result := &CopyResult{
				FileHandle: *job,
			}

			if currentBucket != job.bucketName {
				err := q.toClient.MakeBucket(ctx, job.bucketName, minio.MakeBucketOptions{})
				currentBucket = job.bucketName

				if err != nil {
					result.Err = err
					*results = append(*results, result)
					continue
				}
			}

			reader, err := q.fromClient.readObject(ctx, job)
			if err != nil {
				result.Err = err
				*results = append(*results, result)
				continue
			}

			err = q.toClient.writeObject(ctx, reader, job)
			if err != nil {
				result.Err = err
				*results = append(*results, result)
				continue
			}

		}
	}

	return true
}

func (client *Client) readObject(ctx context.Context, handle *FileHandle) (io.Reader, error) {
	obj, err := client.GetObject(ctx, handle.bucketName, handle.fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (client *Client) writeObject(ctx context.Context, from io.Reader, handle *FileHandle) error {
	_, err := client.PutObject(ctx, handle.bucketName, handle.fileName, from, handle.size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) CopyBuckets(ctx context.Context, to *Client, overwrite bool) ([]*CopyResult, error) {

	files, err := client.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	toFiles, err := to.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*CopyResult, 0, len(files))
	toCopy := make([]*FileHandle, 0, len(files))

	if !overwrite {
		for _, f := range files {
			exists := false

			for _, tf := range toFiles {
				if f.bucketName == tf.bucketName && f.fileName == tf.fileName {
					exists = true
					break
				}
			}

			if exists {
				results = append(results, &CopyResult{
					FileHandle:       *f,
					AlreadyExists:    true,
					DifferentVersion: false,
				})
			} else {
				toCopy = append(toCopy, f)
			}
		}
	} else {
		toCopy = append(toCopy, files...)
	}

	queue := &copyQueue{
		jobs:       make(chan *FileHandle, len(toCopy)),
		fromClient: client,
		toClient:   to,
		ctx:        ctx,
		cancel:     nil, // ???
	}

	queue.toCopy(toCopy)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		queue.doCopy(ctx, &results)
		wg.Done()
	}()

	wg.Wait()

	return results, nil
}
