package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var command = flag.String("command", "", "")
var bucketName = flag.String("bucket", "", "")
var localPath = flag.String("path", "", "")
var endpoint = flag.String("endpoint", "", "")
var disableSSL = flag.Bool("disable-ssl", false, "")

func listObjects(client *s3.S3, bucket *s3.Bucket) {
	log.Printf("bucket: %s\n", aws.StringValue(bucket.Name))
	params := &s3.ListObjectsInput{
		Bucket: bucket.Name,
	}
	objects, err := client.ListObjects(params)
	if err == nil {
		for _, o := range objects.Contents {
			log.Println("object: ", o)
		}
	} else {
		log.Println("Failed to list objects in buckets", err)
	}
}

func listBucket(client *s3.S3) {
	result, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("Failed to list buckets", err)
		return
	}

	log.Println("Buckets:")
	for _, bucket := range result.Buckets {
		log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}
}

func getBucket(client *s3.S3, bucketName string) {
	bucket := &s3.Bucket{Name: &bucketName}
	listObjects(client, bucket)
}

func makePutObjectInput(bucketName string, file *os.File, localPath string) *s3.PutObjectInput {
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	key := path.Base(localPath)
	return &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
}

func putObject(client *s3.S3, bucketName, localPath string) {
	file, err := os.Open(localPath)
	if err != nil {
		log.Printf("err opening file: %s", err)
		return
	}
	defer file.Close()

	result, err := client.PutObject(makePutObjectInput(bucketName, file, localPath))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(result)
}

func putBucket(client *s3.S3, bucket string) {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String("BucketLocationConstraint"),
		},
	}
	result, err := client.CreateBucket(params)

	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(result)
}

func main() {
	flag.Parse()

	// read credentials in shared config file
	shared := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)

	c := aws.NewConfig().
		WithCredentialsChainVerboseErrors(true).
		WithCredentials(shared.Config.Credentials).
		WithRegion(*shared.Config.Region).
		WithEndpoint(*endpoint).
		WithDisableSSL(*disableSSL)
	s, err := session.NewSession(c)
	if err != nil {
		log.Println("Failed to instatiate session", err)
		return
	}

	client := s3.New(s)
	switch *command {
	case "listBucket":
		listBucket(client)
	case "getBucket":
		getBucket(client, *bucketName)
	case "putBucket":
		putBucket(client, *bucketName)
	case "putObject":
		putObject(client, *bucketName, *localPath)
	default:
		log.Println("unknown command", command)
	}
}
