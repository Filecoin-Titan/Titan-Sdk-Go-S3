# titan-sdk-go-s3

The Titan S3 SDK provides a convenient way to interact with the Titan service and the S3-compatible object storage system using Golang. With this SDK, you can easily list buckets, create new buckets, upload objects, and retrieve objects from your Titan S3-compatible storage.

## Installation
To get started with the Titan S3 SDK, you need to have Go installed on your system. Then, you can install the SDK using the following command:

```go
go get github.com/Filecoin-Titan/titan-sdk-go-s3
```

## Usage
Initialize the SDK by obtaining the Titan S3 API client:


```go
address := os.Getenv("LOCATOR_API_INFO")
client := titans3.New(
    config.TitanAddressOption(address),
    config.CandidateIDOption("your-L1-node-id"),
)

s3API, err := client.GetLocalMinioAPI()
if err != nil {
    log.Fatal(err)
}
```
### List Buckets
You can list all the available buckets in your Titan S3-compatible storage:
```go
ctx := context.Background()
buckets, err := s3API.ListBuckets(ctx)
if err != nil {
    log.Fatal(err)
}

for _, bucket := range buckets {
    fmt.Printf("bucket.Name: %s\n", bucket.Name)
}
```

### Create a New Bucket
You can create a new bucket in your Titan S3-compatible storage:
```go
err = s3API.MakeBucket(ctx, "your-bucket-name", minio.MakeBucketOptions{})
if err != nil {
    log.Fatal(err)
}
```

### Upload an Object
You can upload an object to a specific bucket in your Titan S3-compatible storage:
```go
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

```
### Get an Object
You can retrieve an object from a specific bucket in your Titan S3-compatible storage:
```go
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
```

For more examples of how to use the Titan SDK, check out the examples directory in this repository. 

## Issues
Feel free to submit issues and enhancement requests.


## License

See [MIT](LICENSE) for more information.
