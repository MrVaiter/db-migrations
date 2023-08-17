package aws_s3

type FileHandle struct {
	bucketName string
	fileName   string
	size       int64
}

func (f *FileHandle) GetBucketName() string {
	return f.bucketName
}

func (f *FileHandle) GetFileName() string {
	return f.fileName
}

func (f *FileHandle) GetSize() int64 {
	return f.size
}
