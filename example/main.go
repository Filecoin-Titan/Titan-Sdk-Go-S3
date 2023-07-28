package main

import (
	"context"
	"fmt"
	titans3 "github.com/Filecoin-Titan/titan-sdk-go-s3"
	"github.com/Filecoin-Titan/titan-sdk-go-s3/config"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"os"
)

func main() {
	address := os.Getenv("LOCATOR_API_INFO")
	client := titans3.New(
		config.TitanAddressOption(address),
		config.CandidateIDOption("your-L1-node-id"),
	)

	s3API, err := client.GetLocalMinioAPI()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	buckets, err := s3API.ListBuckets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range buckets {
		fmt.Printf("bucket.Name: %s\n", bucket.Name)
	}

	// Creates a new bucket.
	err = s3API.MakeBucket(ctx, "your-bucket-name", minio.MakeBucketOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Upload File Example
	file, err := os.Open("/path/to/your/file.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	uploadInfo, err := s3API.PutObject(ctx, "your-bucket-name", "your-object-name", file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)

	// GetObject Example
	object, err := s3API.GetObject(ctx, "your-bucket-name", "your-object-name", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer object.Close()

	localFile, err := os.Create("/tmp/local-file.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}
}
