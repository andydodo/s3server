package s3disk

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/0x434D53/s3server/common"
)

type Options struct {
	BasePath  string
	CacheSize uint64
}

type Disk struct {
	sync.RWMutex
	Options
}

func (d *Disk) getBucketPath(bucketName string, auth string) string {
	return filepath.Join(d.BasePath, auth, bucketName)
}

func (d *Disk) getFilePath(bucketName string, objectName string, auth string) string {
	path := filepath.Join(d.BasePath, auth, bucketName, objectName)

	return path
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func existsBool(path string) bool {
	b, _ := exists(path)

	return b
}

func (d *Disk) GetService(auth string) (*common.ListAllMyBucketsResult, error) {
	p := d.BasePath

	f, err := os.Open(p)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	fis, err := f.Readdir(0)

	var filtered []os.FileInfo

	for _, fi := range fis {
		if fi.IsDir() {
			filtered = append(filtered, fi)
		}

		fmt.Println(fi.Name())
	}

	res := &common.ListAllMyBucketsResult{Buckets: make([]*common.Bucket, 0)}

	for _, f := range filtered {
		name := f.Name()
		ctime := GetCTime(f)

		res.Buckets = append(res.Buckets, &common.Bucket{Name: name, CreationDate: ctime})
	}

	return res, nil
}

func (d *Disk) DeleteBucket(bucketName string, auth string) error {
	path := d.getBucketPath(bucketName, auth)

	if !existsBool(path) {
		return common.ErrNotFound
	}

	return os.Remove(path)
}

func (d *Disk) GetBucketObjects(bucketName string, auth string) (*common.ListBucketResult, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d *Disk) HeadBucket(bucketName string, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) PutBucket(bucketName string, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) DeleteObject(bucketName string, objectName string, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) GetObject(bucketName string, objectName string, auth string) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (d *Disk) HeadObject(bucketName string, objectName string, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) PutObject(bucketName string, objectName string, data []byte, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) PutObjectCopy(bucketName string, objectName string, targetBucket string, targetObject string, auth string) error {
	return fmt.Errorf("Not implemented")
}

func (d *Disk) PostObject(bucketName string, objectName string, data []byte, auth string) error {
	return fmt.Errorf("Not implemented")
}

func main() {
	var be common.S3Backend

	be = &Disk{}
	fmt.Println(be)
}
