package aws_s3

import (
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestS3(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
