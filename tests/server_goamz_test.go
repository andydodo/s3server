package tests

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

var localRegion = aws.Region{
	Name:                 "test.dev",
	S3Endpoint:           "http://test.dev:10001",
	S3BucketEndpoint:     "",
	S3LocationConstraint: false,
	S3LowercaseBucket:    true,
}

func Reset(t *testing.T) {
	url := localRegion.S3Endpoint + "/_internal/reset"

	_, err := http.Get(url)

	if err != nil {
		t.Fatal("Could not Reset", err)
	}
}

func TestPutGet(t *testing.T) {
	Reset(t)
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

func TestObjectCycle(t *testing.T) {
	Reset(t)
	objectPath := "/test1"
	objectContents := []byte("test1")
	updatedObjectContents := []byte("Updatedtest")

	auth, err := aws.EnvAuth()

	if err != nil {
		t.Error(err)
	}

	s := s3.New(auth, localRegion)

	// Create Bucket
	b := s.Bucket("TestBucket")
	err = b.PutBucket("acl")

	if err != nil {
		log.Fatalf("Couldn't create bucket: %v", err)
	}

	// Put Object
	err = b.Put(objectPath, objectContents, "application/octet-stream", "acl", s3.Options{})

	if err != nil {
		t.Fatal(err)
	}

	// Get the same object
	data, err := b.Get(objectPath)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, objectContents) {
		t.Errorf("Expected content %v, got content: %v", objectContents, data)
	}

	// Check if there is exactly 1 object in the bucket now
	o, err := b.GetBucketContents()

	if err != nil {
		t.Fatalf("Couldn't get bucket contents: %v", err)
	}

	if len(*o) != 1 {
		t.Fatalf("Bucket sould contain 1 object, contains %d objects", len(*o))
	}

	err = b.Put(objectPath, updatedObjectContents, "application/octet-stream", "acl", s3.Options{})

	if err != nil {
		t.Fatalf("Error updateing Object: %v", err)
	}

	// Check that the object has been modified

	data, err = b.Get(objectPath)

	if err != nil {

		t.Fatalf("Error getting updated object: %v", err)
	}

	if !bytes.Equal(data, updatedObjectContents) {
		t.Fatalf("Wrong Bucket contents, expected %v, got %v", data, updatedObjectContents)

	}

	if err != nil {
		t.Fatalf("Couldn't get bucket contents: %v", err)
	}

	if len(*o) != 1 {
		t.Fatalf("Bucket sould contain 1 object after update, contains %d objects", len(*o))
	}

	err = b.Del(objectPath)

	if err != nil {
		t.Fatalf("Error deleting object: %v", err)
	}

	// Check that the bucket is now empty
	o, err = b.GetBucketContents()

	if err != nil {
		t.Fatalf("Couldn't get bucket contents: %v", err)
	}

	if len(*o) != 0 {
		t.Fatalf("Bucket should be empty, but contains %d objects", len(*o))
	}
}

func TestComparedContentMD5AndContent(t *testing.T) {
	t.Fatalf("Not implemented")
}

func TestDeletedOnlyEmptyBuckets(t *testing.T) {
	t.Fatalf("Not implemented")
}
