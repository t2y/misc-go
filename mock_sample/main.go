package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	bucketName = flag.String("bucket", "", "")
	localPath  = flag.String("path", "", "")
	objectKey  = flag.String("objectKey", "", "")
	endpoint   = flag.String("endpoint", "", "")
	disableSSL = flag.Bool("disableSSL", false, "")
)

func writeFile(body io.ReadCloser, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()
	io.Copy(f, body)
	log.Println("wrote file:", fileName)
}

func getObject(client s3iface.S3API, bucketName, key string) {
	bucket := &s3.Bucket{Name: &bucketName}
	params := &s3.GetObjectInput{
		Bucket: bucket.Name,
		Key:    aws.String(key),
	}

	result, err := client.GetObject(params)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer result.Body.Close()
	log.Println("GetObjectOutput:", result)

	fileName := path.Base(key)
	writeFile(result.Body, fileName)
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
		WithEndpoint(*endpoint)

	if *disableSSL {
		c = c.WithDisableSSL(*disableSSL)
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		c = c.WithHTTPClient(client)
	}

	s, err := session.NewSession(c)
	if err != nil {
		log.Println("Failed to instatiate session", err)
		return
	}

	client := s3.New(s)
	getObject(client, *bucketName, *objectKey)
}
