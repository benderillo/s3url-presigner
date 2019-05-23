package svc

import (
	"errors"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage struct {
	s3svc *s3.S3
}

// NewStorage creates a new instance of Storage which can provide presigned URLs for AWS S3 buckets
func NewStorage(awsSession *session.Session) Storage {

	return Storage{
		s3svc: s3.New(awsSession),
	}
}

func (storage Storage) parseS3URI(uri string) (*string, *string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, nil, errors.New("Unable to parse S3 URL: " + uri)
	}

	bucket := parsed.Host
	key := parsed.Path[1:]
	return &bucket, &key, nil
}

func (storage Storage) GetPresignedURL(method string, uri string, expiry time.Duration) (*string, time.Time, error) {
	bucket, key, err := storage.parseS3URI(uri)
	if err != nil {
		return nil, time.Now().UTC(), err
	}

	var req *request.Request
	if method == "put" {
		req, _ = storage.s3svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	} else if method == "get" {
		req, _ = storage.s3svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	}

	if req == nil {
		return nil, time.Now().UTC(), errors.New("Unable to create a request")
	}

	s3url, err := req.Presign(expiry)
	if err != nil {
		return nil, time.Now().UTC(), err
	}

	return &s3url, time.Now().UTC().Add(expiry), nil
}
