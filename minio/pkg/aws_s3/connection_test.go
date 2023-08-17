package aws_s3

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Minio Connection", func() {
	It("Can connect", func() {

		endpoint := "localhost:9000"
		accessKeyID := "ESlDy36Mti8C5JQauVWw"
		secretAccessKey := "Iy66RB8RKUzQZEOnSpi3tS9SdCYVYIgpZHqP2Uss"

		client, err := Connect(context.Background(), endpoint, accessKeyID, secretAccessKey, "")
		Expect(err).To(BeNil())
		Expect(client).NotTo(BeNil())

	})
})
