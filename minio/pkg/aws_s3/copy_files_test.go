package aws_s3

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Copying files", Ordered, func() {
	var from *Client
	var to *Client
	var err error

	BeforeAll(func() {
		from, err = Connect(context.Background(), "localhost:9000", "Dbx4jvuJyRUGSxdLVAmf", "ZoB96j7OoIiEP9TBuCGpUUqdurZFQuxmmpUooDu1", "")
		Expect(err).To(BeNil())
		Expect(from).NotTo(BeNil())

		to, err = Connect(context.Background(), "localhost:9002", "fMj0IHtV4YnrRQpisped", "X1JqBqAuGNE53BHxaeHivXD2spvxDJc3AgqNhtK3", "")
		Expect(err).To(BeNil())
		Expect(to).NotTo(BeNil())

		createdBuckets := from.makeRandomBuckets(5)

		folderPath := "../../../test_files"
		files, err := os.ReadDir(folderPath)
		Expect(err).To(BeNil())
		Expect(files).ToNot(BeNil())

		for _, bucket := range createdBuckets {

			for _, fileName := range files {

				filePath := folderPath + "/" + fileName.Name()
				err = from.uploadFile(context.Background(), bucket, fileName.Name(), filePath)
				Expect(err).To(BeNil())
			}

		}
	})

	It("Can copy files", func() {
		_, err = from.CopyBuckets(context.Background(), to, false)
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		err = from.ClearWithSuffix(context.Background(), "-test-")
		Expect(err).To(BeNil())

		err = to.ClearWithSuffix(context.Background(), "-test-")
		Expect(err).To(BeNil())
	})
})
