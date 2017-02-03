package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var command = flag.String("command", "", "")
var bucketName = flag.String("bucket", "", "")
var localPath = flag.String("path", "", "")
var objectKey = flag.String("objectKey", "", "")
var rangeBytes = flag.String("rangeBytes", "", "")
var endpoint = flag.String("endpoint", "", "")
var disableSSL = flag.Bool("disableSSL", false, "")

func getNextRange(contentRange string, rangeBytes int) string {
	log.Println("contentRange:", contentRange)
	div := strings.Split(contentRange, "/")
	if len(div) == 2 {
		contentLengthBytes, _ := strconv.Atoi(div[1])
		log.Println("contentLengthBytes", contentLengthBytes)
		bytesRange := strings.Split(div[0], "-")
		if len(bytesRange) == 2 {
			lastBytes, err := strconv.Atoi(bytesRange[1])
			if err != nil {
				log.Println(err.Error())
				return ""
			}
			offset := lastBytes + 1
			nextLastBytes := lastBytes + rangeBytes
			if nextLastBytes >= contentLengthBytes {
				nextLastBytes = contentLengthBytes - 1
			}
			nextRange := "bytes=" + strconv.Itoa(offset) + "-" + strconv.Itoa(nextLastBytes)
			log.Println("nextRange:", nextRange)
			return nextRange
		}
	}
	return ""
}

func hasRestFileContents(contentRange string) bool {
	div := strings.Split(contentRange, "/")
	if len(div) == 2 {
		contentLength := div[1]
		bytesRange := strings.Split(div[0], "-")
		if len(bytesRange) == 2 {
			lastByte, _ := strconv.Atoi(bytesRange[1])
			contentLengthByte, _ := strconv.Atoi(contentLength)
			if lastByte+1 != contentLengthByte {
				return true
			}
		}
	}
	return false
}

func writeFile(body io.ReadCloser, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()
	io.Copy(f, body)
	log.Println("wrote file:", fileName)
	log.Println("========================================================================")
}

func concatenateFile(fileName string, maxSegment int) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()

	for i := 0; i <= maxSegment; i++ {
		segmentFileName := getSegmentFileName(fileName, i)
		log.Println("read segment file:", segmentFileName)
		segFile, err := os.Open(segmentFileName)
		if err != nil {
			log.Printf("err opening file: %s", err)
			return
		}
		io.Copy(f, segFile)
	}

	log.Println("complete concatenate from segmented file:", fileName)
}

func getSegmentFileName(fileName string, num int) string {
	return fileName + "." + strconv.Itoa(num)
}

func showProtocolSchema() string {
	if *disableSSL {
		return "Schema: http"
	} else {
		return "Schema: https"
	}
}

func writeFileWithRangeRequest(client *s3.S3, bucketName, key, fileName, contentRange string, rangeBytes int) {
	i := 1
	for {
		params := &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    aws.String(key),
			Range:  aws.String(getNextRange(contentRange, rangeBytes)),
		}
		result, err := client.GetObject(params)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println(showProtocolSchema())
		log.Println("GetObjectOutput:", result)

		writeFile(result.Body, getSegmentFileName(fileName, i))

		if result.ContentRange == nil || !hasRestFileContents(*result.ContentRange) {
			break
		}
		contentRange = *result.ContentRange
		i += 1
	}

	concatenateFile(fileName, i)
}

func getObject(client *s3.S3, bucketName, key, rangeBytes string) {
	bucket := &s3.Bucket{Name: &bucketName}
	params := &s3.GetObjectInput{
		Bucket: bucket.Name,
		Key:    aws.String(key),
	}

	var _rangeBytes int
	if rangeBytes != "" {
		var err error
		_rangeBytes, err = strconv.Atoi(rangeBytes)
		if err != nil {
			log.Println(err.Error())
			return
		}
		rangeString := "bytes=0-" + strconv.Itoa(_rangeBytes-1)
		params = params.SetRange(rangeString)
		log.Println("setRange:", rangeString)
	}

	result, err := client.GetObject(params)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer result.Body.Close()
	log.Println("GetObjectOutput:", result)

	fileName := path.Base(key)
	if result.ContentRange != nil && hasRestFileContents(*result.ContentRange) {
		writeFile(result.Body, fileName+".0")
		writeFileWithRangeRequest(client, bucketName, key, fileName, *result.ContentRange, _rangeBytes)
	} else {
		writeFile(result.Body, fileName)
	}
}

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

func makePutObjectInput(
	bucketName string, file *os.File, localPath string, objectKey string,
) *s3.PutObjectInput {
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	key := objectKey
	if key == "" {
		key = path.Base(localPath)
	}
	return &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
}

func putObject(client *s3.S3, bucketName, localPath string, objectKey string) {
	file, err := os.Open(localPath)
	if err != nil {
		log.Printf("err opening file: %s", err)
		return
	}
	defer file.Close()

	params := makePutObjectInput(bucketName, file, localPath, objectKey)
	log.Println("PutObjectInput:", params)
	result, err := client.PutObject(params)
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
	switch *command {
	case "listBucket":
		listBucket(client)
	case "getBucket":
		getBucket(client, *bucketName)
	case "getObject":
		getObject(client, *bucketName, *objectKey, *rangeBytes)
	case "putBucket":
		putBucket(client, *bucketName)
	case "putObject":
		putObject(client, *bucketName, *localPath, *objectKey)
	default:
		log.Println("unknown command", command)
	}
}
