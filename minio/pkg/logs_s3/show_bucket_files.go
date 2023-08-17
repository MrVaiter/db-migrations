package logs

import (
	. "example.com/pkg/aws_s3"
	"fmt"
)

func ShowBucketFiles(files []*FileHandle) {
	fmt.Println("File\t\t|Size")

	bucketName := ""
	
	for _, file := range files {
		if bucketName != file.GetBucketName(){
			fmt.Println("---------------------")
			fmt.Println("Bucket: ",file.GetBucketName())
			fmt.Println("---------------------")
		}

		fmt.Println(file.GetFileName(), "\t|", file.GetSize(), "b")
		bucketName = file.GetBucketName()
	}
}
