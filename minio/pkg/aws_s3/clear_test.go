package aws_s3

import (
	"context"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deleting buckets", func() {
	It("Can delete buckets", func() {

		endpoint, ok := os.LookupEnv("S3_ENDPOINT")
		Expect(ok).To(BeTrue())

		accessKeyID, ok := os.LookupEnv("S3_ACCESS_KEY_ID")
		Expect(ok).To(BeTrue())

		secretAccessKey, ok := os.LookupEnv("S3_SECRET_KEY_ID")
		Expect(ok).To(BeTrue())

		client, err := Connect(context.Background(), endpoint, accessKeyID, secretAccessKey, "")
		Expect(err).To(BeNil())
		Expect(client).NotTo(BeNil())

		client.makeRandomBuckets(5)

		err = client.ClearWithSuffix("-test-")
		Expect(err).To(BeNil())
		
		buckets, err := client.ListBuckets(context.Background())
		Expect(err).To(BeNil())

		for _, bucket := range buckets {
			Expect(strings.Contains(bucket.Name, "-test-")).To(BeFalse())
		}
	})
})
