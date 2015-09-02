package s3backend

import "testing"

var backend S3Backend

var create func() S3Backend

func init() {
	create = NewS3InMemory
}

func TestPutHeadBucket(t *testing.T) {
	backend := create()

	err := backend.PutBucket("TestBucket", "")

	if err != nil {
		t.Error("PutBucket gave an Error")
	}

	err = backend.HeadBucket("TestBucket", "")

	if err != nil {
		t.Error("Head for existing Bucket resulted in an error")
	}

}
func TestHeadNonExistBucket(t *testing.T) {
	backend := create()

	err := backend.HeadBucket("TestBucket", "")

	if err != ErrNotFound {
		t.Errorf("Head for non existing Bucket should give an ErrNotFound")
	}
}
