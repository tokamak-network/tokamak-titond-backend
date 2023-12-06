package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	output *s3manager.UploadOutput
	err    error
}

func (mock *MockUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return mock.output, mock.err
}

func TestUploadFile(t *testing.T) {

	s3 := NewS3(&S3Config{
		BucketName: "test bucket",
	})
	mockUploader := &MockUploader{}
	mockUploader.output = &s3manager.UploadOutput{
		Location: "localtion",
	}
	mockUploader.err = nil
	s3.uploader = mockUploader

	result, err := s3.UploadContent("filename", "content")
	fmt.Println(result, err)
	assert.Equal(t, nil, err)
}

func TestUploadFileFailed(t *testing.T) {

	s3 := NewS3(&S3Config{
		BucketName: "test bucket",
	})
	mockUploader := &MockUploader{}
	mockUploader.err = errors.New("")
	s3.uploader = mockUploader

	result, err := s3.UploadContent("filename", "content")
	fmt.Println(result, err)
	assert.NotEqual(t, nil, err)
}
