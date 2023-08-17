package aws_s3

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/minio/minio-go/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Listing Connection", Ordered, func() {

	var client *Client
	var bucketName string = "eleven-test"
	var objectName string = "11.png"
	var filePath string = "../../../test_files/10.png"

	BeforeAll(func() {
		client, bucketName = setUpEnvironment(bucketName, objectName, filePath)
	})

	It("Can list buckets", func() {

		buckets, err := client.ListBuckets(context.Background())
		Expect(err).To(BeNil())
		Expect(buckets).NotTo(BeEmpty())

	})

	It("Can list objects", func() {

		files, err := client.ListFiles(context.Background())
		Expect(err).To(BeNil())
		Expect(files).NotTo(BeEmpty())

	})

	AfterAll(func() {
		err := clear(client)
		Expect(err).To(BeNil())
	})
})

func setUpEnvironment(bucketName string, objectName string, filePath string) (*Client, string) {
	endpoint, ok := os.LookupEnv("S3_ENDPOINT")
	Expect(ok).To(BeTrue())

	accessKeyID, ok := os.LookupEnv("S3_ACCESS_KEY_ID")
	Expect(ok).To(BeTrue())

	secretAccessKey, ok := os.LookupEnv("S3_SECRET_KEY_ID")
	Expect(ok).To(BeTrue())

	client, err := Connect(context.Background(), endpoint, accessKeyID, secretAccessKey, "")
	Expect(err).To(BeNil())

	bucketName = fmt.Sprintf("%v-test-%v", bucketName, rand.Intn(100))

	// Create test bucket
	err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	Expect(err).To(BeNil())

	err = uploadFile(client, bucketName, objectName, filePath)
	Expect(err).To(BeNil())

	return client, bucketName
}
