package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	s := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
	s3Client := s3.New(s)
	result, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("Failed to list buckets", err)
		return
	}

	log.Println("Buckets:")
	for _, bucket := range result.Buckets {
		log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
		objects, err := s3Client.ListObjects(&s3.ListObjectsInput{Bucket: bucket.Name})
		if err == nil {
			for _, o := range objects.Contents {
				log.Println("object: ", o)
			}
		} else {
			log.Println("Failed to list objects in buckets", err)
		}

	}
}
