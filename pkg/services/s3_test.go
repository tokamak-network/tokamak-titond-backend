package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

var S3DisableContentMD5Validation bool = true

var sess = func() *session.Session {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(true),
		Endpoint:   aws.String(server.URL),
		Region:     aws.String("Test"),
	}))
}()

func TestUploadFile(t *testing.T) {

	s3 := NewS3(&S3Config{
		BucketName: "test bucket",
	})
	s3.sess = sess

	result, err := s3.UploadContent("filename", "content")
	fmt.Println(result, err)
	assert.Equal(t, nil, err)
}

func TestUploadFileFailed(t *testing.T) {

	s3 := NewS3(&S3Config{})
	s3.sess = sess

	result, err := s3.UploadContent("filename", "content")
	fmt.Println(result, err)
	assert.NotEqual(t, nil, err)
}
