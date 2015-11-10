package main

import (
	"bytes"
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

func TestBucketCycle(t *testing.T) {
	auth, err := aws.EnvAuth()

	if err != nil {
		t.Error(err)
	}

	s := s3.New(auth, localRegion)

	b := s.Bucket("TestBucket")

	err = b.PutBucket("acl")

	if err != nil {
		t.Fatal(err)
	}

	o, err := b.GetBucketContents()

	if err != nil {
		t.Fatal(err)
	}

	if len(*o) != 0 {
		t.Fatalf("Bucket should be empty, but has %d object", len(*o))
	}
}

func _TestGetBucket(t *testing.T) {

}

func byteSliceCompare() {

}

func TestObjectCycle(t *testing.T) {
	objectPath := "/test1"
	objectContents := []byte("test1")

	auth, err := aws.EnvAuth()

	if err != nil {
		t.Error(err)
	}

	s := s3.New(auth, localRegion)

	b := s.Bucket("TestBucket")
	err = b.Put(objectPath, objectContents, "application/octet-stream", "acl", s3.Options{})

	if err != nil {
		t.Error(err)
	}

	data, err := b.Get(objectPath)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, objectContents) {
		t.Errorf("Expected content %v, got content: %v", objectContents, data)
	}
}
