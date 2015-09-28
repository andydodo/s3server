package main

import (
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

var localRegion = aws.Region{
	Name:                 "localhost",
	S3Endpoint:           "http://localhost:10001",
	S3BucketEndpoint:     "",
	S3LocationConstraint: false,
	S3LowercaseBucket:    true,
}

func TestPutBucket(t *testing.T) {
	auth, err := aws.EnvAuth()

	if err != nil {
		t.Error(err)
	}

	s := s3.New(auth, localRegion)

	b := s.Bucket("TestBucket")

	err = b.PutBucket("acl")

	if err != nil {
		t.Error(err)
	}
}

func TestGetBucket(t *testing.T) {

}

func TestPutObject(t *testing.T) {
	auth, err := aws.EnvAuth()

	if err != nil {
		t.Error(err)
	}

	s := s3.New(auth, localRegion)

	b := s.Bucket("TestBucket")
	err = b.Put("/test1", []byte("test1"), "application/octet-stream", "acl", s3.Options{})

	if err != nil {
		t.Error(err)
	}
}
