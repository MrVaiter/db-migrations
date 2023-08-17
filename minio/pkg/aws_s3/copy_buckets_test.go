package aws_s3

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/minio/minio-go/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Copying buckets", Ordered, func() {
	var from *Client
	var to *Client
	var err error
	var copiedBuckets []string

	BeforeAll(func() {
		from, err = Connect(context.Background(), "localhost:9000", "Dbx4jvuJyRUGSxdLVAmf", "ZoB96j7OoIiEP9TBuCGpUUqdurZFQuxmmpUooDu1", "")
		Expect(err).To(BeNil())
		Expect(from).NotTo(BeNil())

		to, err = Connect(context.Background(), "localhost:9002", "fMj0IHtV4YnrRQpisped", "X1JqBqAuGNE53BHxaeHivXD2spvxDJc3AgqNhtK3", "")
		Expect(err).To(BeNil())
		Expect(to).NotTo(BeNil())

		copiedBuckets = makeRandomBuckets(from, 10)
	})

	It("Can copy buckets", func() {
		err = CopyBuckets(from, to, "-test-")
		Expect(err).To(BeNil())

		for _, bucket := range copiedBuckets {
			result, err := to.BucketExists(context.Background(), bucket)
			Expect(result).To(BeTrue())
			Expect(err).To(BeNil())
		}
	})

	AfterAll(func() {
		err = clear(from)
		Expect(err).To(BeNil())

		err = clear(to)
		Expect(err).To(BeNil())
	})
})

func makeRandomBuckets(client *Client, amount int) []string {

	if amount <= 0 {
		return nil
	}

	names := make([]string, 0, 10)
	ctx := context.Background()

	for i := 1; i <= amount; i++ {

		bucketName := fmt.Sprintf("bucket-test-%v", rand.Intn(10000))
		names = append(names, bucketName)
		err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		Expect(err).To(BeNil())
	}

	return names
}
