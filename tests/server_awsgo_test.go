package tests

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestAWSFirstRequest(t *testing.T) {
	Reset(t)
	c := &aws.Config{Endpoint: aws.String("http://test.dev:10001"), Region: aws.String("us-west-2")}

	svc := s3.New(session.New(c))

	bucket := "test1"

	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})

	if err != nil {
		t.Fatal(err)
	}
}

func TestAWSObjectCycle(t *testing.T) {

}
