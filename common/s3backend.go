package common

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("Not found")
	ErrForbidden     = errors.New("Forbidden")
	ErrAlreadyExists = errors.New("Already Exists")
)

type ListAllMyBucketsResult struct {
	Owner   Owner
	Buckets []*Bucket
}

type Bucket struct {
	Name         string
	CreationDate time.Time
}

type ListBucketResult struct {
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	IsTruncated bool
	Contents    []Contents
}

type Contents struct {
	Key          string
	LastModified *time.Time
	ETag         string
	Size         int
	StorageClass string
	Owner        Owner
}

type Metadata struct{}

type S3Backend interface {
	GetService(auth string) (*ListAllMyBucketsResult, error)
	DeleteBucket(bucket string, auth string) error
	GetBucketObjects(bucket string, auth string) (*ListBucketResult, error)
	HeadBucket(bucket string, auth string) error
	PutBucket(bucket string, auth string) error // More Parameters available
	DeleteObject(bucket string, object string, auth string) error
	GetObject(bucket string, object string, auth string) ([]byte, string, error)
	//GetObjectStream(bucket string, object string, auth string) (io.WriteCloser, error)
	HeadObject(bucket string, object string, auth string) error
	PutObject(bucket string, object string, data []byte, contentType string, auth string) error
	//PutObjectStream(bucket string, object string, r io.ReadCloser, auth string) error
	PutObjectCopy(bucket string, object string, targetBucket string, targetObject string, auth string) error
	PostObject(bucket string, object string, data []byte, contentType string, auth string) error
	//PostObjectStream(bucket string, object string, r io.ReadCloser, auth string) error
	Reset()
}

type Options struct {
	SSE              bool
	Meta             map[string][]string
	ContentEncoding  string
	CacheControl     string
	RedirectLocation string
	ContentMD5       string
	// x-amz-storage-class []string
}

type CopyOptions struct {
	Options
	MetadataDirective string
	ContentType       string
}

// CopyObjectResult is the output from a Copy request
type CopyObjectResult struct {
	ETag         string
	LastModified string
}

type ACL string

const (
	Private           = ACL("private")
	PublicRead        = ACL("public-read")
	PublicReadWrite   = ACL("public-read-write")
	AuthenticatedRead = ACL("authenticated-read")
	BucketOwnerRead   = ACL("bucket-owner-read")
	BucketOwnerFull   = ACL("bucket-owner-full-control")
)

type Delete struct {
	Quiet   bool     `xml:"Quiet,omitempty"`
	Objects []Object `xml:"Object"`
}

type Object struct {
	Key       string `xml:"Key"`
	VersionId string `xml:"VersionId,omitempty"`
}

type ListResp struct {
	Name       string
	Prefix     string
	Delimiter  string
	Marker     string
	NextMarker string
	MaxKeys    int
	IsTruncated    bool
	Contents       []Key
	CommonPrefixes []string `xml:">Prefix"`
}

// The Key type represents an item stored in an S3 bucket.
type Key struct {
	Key          string
	LastModified string
	Size         int64
	// ETag gives the hex-encoded MD5 sum of the contents,
	// surrounded with double-quotes.
	ETag         string
	StorageClass string
	Owner        Owner
}

type VersionsResp struct {
	Name            string
	Prefix          string
	KeyMarker       string
	VersionIdMarker string
	MaxKeys         int
	Delimiter       string
	IsTruncated     bool
	Versions        []Version
	CommonPrefixes  []string `xml:">Prefix"`
}

type Version struct {
	Key          string
	VersionId    string
	IsLatest     bool
	LastModified string
	ETag         string
	Size         int64
	Owner        Owner
	StorageClass string
}

type Error struct {
	StatusCode int    // HTTP status code (200, 403, ...)
	Code       string // EC2 error code ("UnsupportedOperation", ...)
	Message    string // The human-oriented error message
	Resource string
	RequestId  string
}

type Owner struct {
	ID          string
	DisplayName string
}

type RequestHeaders struct {
	Authorization      string
	ContentLength      uint64
	ContentType        string
	ContentMD5         string
	Date               time.Time
	Expect             string // can be 100-continue
	Host               string
	XAmzContentSha256  string
	XAmzDate           time.Time
	XAmzSecurityTokens string
}

type ResponseHeaders struct {
	ContentLength    string
	ContentType      string
	Connection       string // open | close
	Date             string
	ETag             string
	Server           string
	XAmzDeleteMarker bool
	XAmzId2          string
	XAmzRequestId    string
	XAmzVersionId    string
}
