package services

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IFIleManager interface {
	UploadContent(fileName string, content string) (string, error)
}

type S3 struct {
	bucket string
	sess   *session.Session
}

func (s3 *S3) initialize() {
	s3.sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")}))
}

func (s3 *S3) UploadContent(fileName string, content string) (string, error) {
	fmt.Printf(" Address: %p \n", s3)
	uploader := s3manager.NewUploader(s3.sess)

	var buffer bytes.Buffer

	_, err := buffer.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to buffer:", err)
		return "", err
	}

	reader := bytes.NewReader(buffer.Bytes())

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s3.bucket,
		Key:    &fileName,
		Body:   reader,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("File is uploaded to, %s\n", result.Location)
	return result.Location, nil
}

func newS3() *S3 {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "titond"
	}
	fmt.Println("bucket", bucket)
	s3 := &S3{bucket: bucket}
	s3.initialize()
	return s3
}

var fileManager IFIleManager = newS3()

func S3Service() IFIleManager {
	return fileManager
}
