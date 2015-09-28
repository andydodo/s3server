package s3backend

import (
	"sync"

	. "github.com/0x434D53/s3server/common"
)

type object struct {
	name     string
	contents []byte
}

type bucket struct {
	objects map[string]*object
	sync.Mutex
}

type S3InMemory struct {
	buckets map[string]*bucket
	sync.Mutex
}

func NewS3InMemory() S3Backend {
	s3 := &S3InMemory{
		buckets: make(map[string]*bucket, 0),
	}

	return s3
}

func (s3 *S3InMemory) PostObject(bucketName string, objectName string, data []byte, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	b, ok := s3.buckets[bucketName]

	if !ok {
		return ErrNotFound
	}

	b.objects[objectName] = &object{name: objectName, contents: data}

	return nil
}

func (s3 *S3InMemory) PutObjectCopy(bucketName string, objectName string, targetBucketName string, targetObjectName string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	b, err := s3.GetObject(bucketName, objectName, auth)

	if err != nil {
		return err
	}

	err = s3.PutObject(targetBucketName, targetObjectName, b, auth)

	if err != nil {
		return err
	}

	return nil
}

func (s3 *S3InMemory) PutObject(bucketName string, objectName string, data []byte, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	b, ok := s3.buckets[bucketName]
	if !ok {
		return ErrNotFound
	}

	b.objects[objectName] = &object{name: objectName, contents: data}

	return nil
}

func (s3 *S3InMemory) DeleteBucket(bucketName string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	delete(s3.buckets, bucketName)

	return nil
}

func (s3 *S3InMemory) PutBucket(bucketName string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	if _, ok := s3.buckets[bucketName]; !ok {
		s3.buckets[bucketName] = &bucket{objects: make(map[string]*object, 0)}
		return nil
	}

	return ErrAlreadyExists
}

func (s3 *S3InMemory) GetBucketObjects(bucketName string, auth string) ([]string, error) {
	b, ok := s3.buckets[bucketName]

	if !ok {
		return nil, ErrNotFound
	}

	b.Lock()
	defer b.Unlock()

	ret := make([]string, 0, len(b.objects))

	for k, _ := range b.objects {
		ret = append(ret, k)
	}

	return ret, nil
}

func (s3 *S3InMemory) HeadBucket(bucket string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	if _, ok := s3.buckets[bucket]; !ok {
		return ErrNotFound
	}

	return nil
}

func (s3 *S3InMemory) HeadObject(bucket string, object string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	if b, ok := s3.buckets[bucket]; !ok {
		return ErrNotFound
	} else {
		if _, ok := b.objects[object]; !ok {
			return ErrNotFound
		}
	}

	return nil
}

func (s3 *S3InMemory) DeleteObject(bucket string, object string, auth string) error {
	s3.Lock()
	defer s3.Unlock()

	return nil
}

func (s3 *S3InMemory) GetObject(bucket string, object string, auth string) ([]byte, error) {
	s3.Lock()
	defer s3.Unlock()

	b, ok := s3.buckets[bucket]

	if !ok {
		return nil, ErrNotFound
	}

	o, ok := b.objects[object]

	if !ok {
		return nil, ErrNotFound
	}
	return o.contents, nil
}
