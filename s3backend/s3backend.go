package s3backend

import "errors"

var (
	ErrNotFound      = errors.New("Not found")
	ErrForbidden     = errors.New("Forbidden")
	ErrAlreadyExists = errors.New("Already Exists")
)

type Metadata struct{}

type S3Backend interface {
	GetService(auth string) error
	DeleteBucket(bucket string, auth string) error
	GetBucketObjects(bucket string, auth string) ([]string, error)
	HeadBucket(bucket string, auth string) error
	PutBucket(bucket string, auth string) error // More Parameters available
	DeleteObject(bucket string, object string, auth string) error
	GetObject(bucket string, object string, auth string) ([]byte, error)
	//GetObjectStream(bucket string, object string, auth string) (io.WriteCloser, error)
	HeadObject(bucket string, object string, auth string) error
	PutObject(bucket string, object string, data []byte, auth string) error
	//PutObjectStream(bucket string, object string, r io.ReadCloser, auth string) error
	PutObjectCopy(bucket string, object string, targetBucket string, targetObject string, auth string) error
	PostObject(bucket string, object string, data []byte, auth string) error
	//PostObjectStream(bucket string, object string, r io.ReadCloser, auth string) error
}
