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

func (from *Client) CopyFiles(ctx context.Context, to *Client, filter string) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := from.CopyBucketsWithSuffix(context.Background(), to, filter)
	if err != nil {
		return err
	}

	log.Print("Starting copying files from ", from.EndpointURL().Host, " to ", to.EndpointURL().Host)
	fromFiles, err := from.ListFiles(ctx)
	if err != nil {
		return err
	}

	toFiles, err := to.ListFiles(ctx)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, file := range fromFiles {

		if strings.Contains(file.bucketName, filter) {

			buffer.WriteString(file.bucketName + "/" + file.fileName)
			if !isContain(toFiles, file) {
				reader, err := from.readObject(ctx, file)
				if err != nil {
					return err
				}

				_, err = to.PutObject(ctx,
					file.bucketName,
					file.fileName,
					reader,
					file.size,
					minio.PutObjectOptions{})
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

func isContain(list []*FileHandle, object *FileHandle) bool {
	for _, file := range list {
		if file.bucketName == object.bucketName && file.fileName == object.fileName {
			return true
		}
	}

	return false
}

func (client *Client) uploadFile(ctx context.Context, bucketName string, fileName string, filePath string) error {
	// Upload test object
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx,
		bucketName,
		fileName,
		file,
		fileInfo.Size(),
		minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
