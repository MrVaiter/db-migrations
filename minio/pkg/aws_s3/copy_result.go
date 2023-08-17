package aws_s3

type CopyResult struct {
	FileHandle
	AlreadyExists    bool
	DifferentVersion bool
	Err              error
}
