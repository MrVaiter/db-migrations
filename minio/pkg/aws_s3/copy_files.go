package aws_s3

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (from *Client) CopyFiles(to *Client, filter string) error {
	ctx := context.Background()

	err := from.CopyBucketsWithSuffix(to, filter)
	if err != nil {
		return err
	}

	fmt.Println("\nStarting copying files from", from.EndpointURL().Host, "to", to.EndpointURL().Host)
	fromFiles, err := from.ListFiles(ctx)
	if err != nil {
		return err
	}

	toFiles, err := to.ListFiles(ctx)
	if err != nil {
		return err
	}

	for _, file := range fromFiles {

		if strings.Contains(file.bucketName, filter) {

			fmt.Print(file.bucketName, "/", file.fileName)
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

func isContain(list []*FileHandle, object *FileHandle) bool {
	for _, file := range list {
		if file.bucketName == object.bucketName && file.fileName == object.fileName {
			return true
		}
	}

	return false
}

func (client *Client) uploadFile(bucketName string, fileName string, filePath string) error {
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

	_, err = client.PutObject(context.Background(),
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
