package main

import (
	"context"
	"log"
	"os"

	. "example.com/pkg/aws_s3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	endpoint, ok := os.LookupEnv("S3_ENDPOINT")
	if !ok {
		panic("S3_ENDPOINT missing!")
	}
	accessKeyID, ok := os.LookupEnv("S3_ACCESS_KEY_ID")
	if !ok {
		panic("S3_ACCESS_KEY_ID missing!")
	}
	secretAccessKey, ok := os.LookupEnv("S3_SECRET_KEY_ID")
	if !ok {
		panic("S3_SECRET_KEY_ID missing!")
	}

	client, err := Connect(context.Background(), endpoint, accessKeyID, secretAccessKey, "")
	if err != nil {
		log.Fatalln(err)
	}

	another_client, err := Connect(context.Background(),
		"localhost:9002",
		"RFhwTT46QbsUlnCT3BYV",
		"tgshup5CoKMJHxjlIUPsIRIun50sqR2ZI3lrN6Hd",
		"")

	// TODO: delete
	another_client.ClearWithSuffix(context.Background(), "")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = client.CopyBuckets(context.Background(), another_client, false)
	if err != nil {
		log.Fatalln(err)
	}

	// Clear(another_client)
}
