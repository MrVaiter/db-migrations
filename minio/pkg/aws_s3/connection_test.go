package aws_s3

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Minio Connection", func() {
	It("Can connect", func() {

		endpoint, ok := os.LookupEnv("S3_ENDPOINT")
		Expect(ok).ToNot(BeNil())

		accessKeyID, ok := os.LookupEnv("S3_ACCESS_KEY_ID")
		Expect(ok).ToNot(BeNil())

		secretAccessKey, ok := os.LookupEnv("S3_SECRET_KEY_ID")
		Expect(ok).ToNot(BeNil())

		client, err := Connect(context.Background(), endpoint, accessKeyID, secretAccessKey, "")
		Expect(err).To(BeNil())
		Expect(client).NotTo(BeNil())

	})
})
