package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ResetAWS(t *testing.T) {
	url := "http://test.dev:10001/_internal/reset"

	_, err := http.Get(url)

	if err != nil {
		t.Fatal("Could not Reset", err)
	}
}

func initTest(t *testing.T) *s3.S3 {
	ResetAWS(t)

	c := &aws.Config{Endpoint: aws.String("http://test.dev:10001"), Region: aws.String("us-west-2")}
	svc := s3.New(session.New(c))
	return svc
}

func TestAWSFirstRequest(t *testing.T) {
	svc := initTest(t)
	bucket := "test1"

	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})

	if err != nil {
		t.Fatal(err)
	}
}

func TestAWSPutGet(t *testing.T) {
	svc := initTest(t)

	bucket := "testBucket"

	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})

	if err != nil {
		t.Fatal(err)
	}

	loo, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		t.Fatal(err)
	}

	if len(loo.Contents) != 0 {
		t.Fatalf("Bucket should be empty, but has %d object", len(loo.Contents))
	}
}

func TestAWSObjectCyclePathMethod(t *testing.T) {
	bucket := "TestBucket"
	awsObjectCycle(bucket, t)
}

func testAWSObjectCycleSubdomainMethod(t *testing.T) {
	bucket := "testbucket"
	awsObjectCycle(bucket, t)
}

func awsObjectCycle(bucket string, t *testing.T) {

	objectPath := "test1"
	objectContents := []byte("test1")
	updatedObjectContents := []byte("Updatedtest")

	svc := initTest(t)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})

	if err != nil {
		t.Fatal(err)
	}

	// Put Object
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectPath),
		Body:   bytes.NewReader(objectContents),
	})

	if err != nil {
		t.Fatal(err)
	}

	// Get the same object
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectPath),
	})

	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	if !bytes.Equal(data, objectContents) {
		t.Errorf("Expected content %v, got content: %v", objectContents, data)
	}

	// Check if there is exactly 1 object in the bucket now
	loo, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		t.Fatalf("Couldn't get bucket contents: %v", err)
	}

	if len(loo.Contents) != 1 {
		t.Fatalf("Bucket sould contain 1 object, contains %d objects", len(loo.Contents))
	}

	// Update the object
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectPath),
		Body:   bytes.NewReader(updatedObjectContents),
	})

	if err != nil {
		t.Fatal(err)
	}

	// Check if the object has been modified
	resp, err = svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectPath),
	})

	if err != nil {
		t.Fatal(err)
	}

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	if !bytes.Equal(data, updatedObjectContents) {
		t.Errorf("Expected content %s, got content: %s", string(objectContents), string(data))
	}

	// Check that there is still exactly one object in the bucket
	loo, err = svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		t.Fatalf("Couldn't get bucket contents: %v", err)
	}

	if len(loo.Contents) != 1 {
		t.Fatalf("Bucket sould contain 1 object after update, contains %d objects", len(loo.Contents))
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectPath),
	})

	if err != nil {
		t.Fatalf("Error deleting object: %v", err)
	}

	// Check that the bucket is now empty
	loo, err = svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		t.Fatalf("Couldn't get bucket after delete  contents: %v", err)
	}

	if len(loo.Contents) != 0 {
		t.Fatalf("Bucket should be empty after delete, contains %d objects", len(loo.Contents))
	}
}
