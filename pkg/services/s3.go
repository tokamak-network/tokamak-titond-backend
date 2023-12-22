package services

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IFIleManager interface {
	UploadContent(fileName string, content string) (string, error)
}

type S3Config struct {
	BucketName string
	AWSRegion  string
}

type IUploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type S3 struct {
	config   *S3Config
	sess     *session.Session
	uploader IUploader
}

func (s3 *S3) initialize() {
	fmt.Println("S3 config:", s3.config)
	s3.sess = session.Must(session.NewSession(&aws.Config{
		Region: &s3.config.AWSRegion}))
	s3.uploader = s3manager.NewUploader(s3.sess)
}

func (s3 *S3) UploadContent(fileName string, content string) (string, error) {

	var buffer bytes.Buffer

	buffer.WriteString(content)

	reader := bytes.NewReader(buffer.Bytes())

	result, err := s3.uploader.Upload(&s3manager.UploadInput{
		Bucket:      &s3.config.BucketName,
		Key:         &fileName,
		Body:        reader,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("File is uploaded to, %s\n", result.Location)
	return result.Location, nil
}

func NewS3(config *S3Config) *S3 {
	s3 := &S3{config: config}
	s3.initialize()
	return s3
}
