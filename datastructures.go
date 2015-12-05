package main

import (
	"fmt"
	"time"
	"github.com/0x434D53/s3server/common"
)

type S3METHOD int

func (s S3METHOD) String() string {
	switch s {
	case GETBUCKET_OBJECTLIST:
		return "GET Bucket objectlist"
	case GETBUCKET:
		return "GET Bucket"
	case GETBUCKET_ACL:
		return "GET Bucket acl"
	case GETBUCKET_CORS:
		return "GET Bucket cors"
	case GETBUCKET_LIFECYCLE:
		return "GET Bucket lifecycle"
	case GETBUCKET_POLICY:
		return "GET Bucket policy"
	case GETBUCKET_LOCATION:
		return "GET Bucket location"
	case GETBUCKET_LOGGING:
		return "GET Bucket logging"
	case GETBUCKET_NOTIFICATION:
		return "GET Bucket notification"
	case GETBUCKET_REPLICATION:
		return "GET Bucket replication"
	case GETBUCKET_TAGGING:
		return "GET Bucket tagging"
	case GETBUCKET_OBJECTVERSION:
		return "GET Bucket objectversion"
	case GETBUCKET_REQUESTPAYMENT:
		return "GET Bucket requestPayment"
	case GETBUCKET_VERSIONING:
		return "GET Bucket versioning"
	case GETBUCKET_WEBSITE:
		return "GET Bucket website"
	case DELETEBUCKET:
		return "DELETE Bucket"
	case DELETEBUCKET_CORS:
		return "DELETE Bucket cors"
	case DELETEBUCKET_LIFTCYCLE:
		return "DELETE Bucket lifecycle"
	case DELETEBUCKET_POLICY:
		return "DELETE Bucket policy"
	case DELETEBUCKET_REPLICATION:
		return "DELETE Bucket replication"
	case DELETEBUCKET_TAGGING:
		return "DELETE Bucket taggin"
	case DELETEBUCKET_WEBSITE:
		return "DELETE Bucket website"
	case PUTBUCKET:
		return "PUT Bucket"
	case PUTBUCKET_ACL:
		return "PUT Bucket acl"
	case PUTBUCKET_CORS:
		return "PUT Bucket cors"
	case PUTBUCKET_LIFECYCLE:
		return "PUT Bucket lifecycle"
	case PUTBUCKET_POLICY:
		return "PUT Bucket policy"
	case PUTBUCKET_LOGGING:
		return "PUT Bucket logging"
	case PUTBUCKET_NOTIFICATION:
		return "PUT Bucket notification"
	case PUTBUCKET_REPLICATION:
		return "PUT Bucket replication"
	case PUTBUCKET_TAGGING:
		return "PUT Bucket tagging"
	case PUTBUCKET_REQUESTPAYMENT:
		return "PUT Bucket requestPayment"
	case PUTBUCKET_VERSIONING:
		return "PUT Bucket versioning"
	case PUTBUCKET_WEBSITE:
		return "PUT Bucket website"
	case HEADBUCKET:
		return "HEAD Bucket"
	case DELETEOBJECT:
		return "DELETE Object"
	case GETOBJECT:
		return "GET Object"
	case GETOBJECT_ACL:
		return "GET Object acl"
	case GETOBJECT_TORRENT:
		return "GET Object torrent"
	case HEADOBJECT:
		return "HEAD Object"
	case POSTOBJECT:
		return "POST Object"
	case POSTOBJECT_RESTORE:
		return "POST Object restore"
	case PUTOBJECT:
		return "PUT Object"
	case PUTOBJECT_ACL:
		return "PUT Object acl"
	}

	return ""
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

type Delete struct {
	Quiet   bool     `xml:"Quiet,omitempty"`
	Objects []Object `xml:"Object"`
}

type Object struct {
	Key       string `xml:"Key"`
	VersionId string `xml:"VersionId,omitempty"`
}

// The Owner type represents the owner of the object in an S3 bucket.
type Owner struct {
	ID          string
	DisplayName string
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

type Options struct {
	SSE              bool
	Meta             map[string][]string
	ContentEncoding  string
	CacheControl     string
	RedirectLocation string
	ContentMD5       string
	// What else?
	// Content-Disposition string
	//// The following become headers so they are []strings rather than strings... I think
	// x-amz-storage-class []string
}

// The VersionsResp type holds the results of a list bucket Versions operation.
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

// The Version type represents an object version stored in an S3 bucket.
type Version struct {
	Key          string
	VersionId    string
	IsLatest     bool
	LastModified string
	// ETag gives the hex-encoded MD5 sum of the contents,
	// surrounded with double-quotes.
	ETag         string
	Size         int64
	Owner        Owner
	StorageClass string
}

type Error struct {
	StatusCode int    // HTTP status code (200, 403, ...)
	Code       string // EC2 error code ("UnsupportedOperation", ...)
	Message    string // The human-oriented error message
	BucketName string
	RequestId  string
	HostId     string
}

const (
	GETBUCKET_OBJECTLIST S3METHOD = iota
	GETBUCKET
	GETBUCKET_ACL
	GETBUCKET_CORS
	GETBUCKET_LIFECYCLE
	GETBUCKET_POLICY
	GETBUCKET_LOCATION
	GETBUCKET_LOGGING
	GETBUCKET_NOTIFICATION
	GETBUCKET_REPLICATION
	GETBUCKET_TAGGING
	GETBUCKET_OBJECTVERSION
	GETBUCKET_REQUESTPAYMENT
	GETBUCKET_VERSIONING
	GETBUCKET_WEBSITE
	DELETEBUCKET
	DELETEBUCKET_CORS
	DELETEBUCKET_LIFTCYCLE
	DELETEBUCKET_POLICY
	DELETEBUCKET_REPLICATION
	DELETEBUCKET_TAGGING
	DELETEBUCKET_WEBSITE
	PUTBUCKET
	PUTBUCKET_ACL
	PUTBUCKET_CORS
	PUTBUCKET_LIFECYCLE
	PUTBUCKET_POLICY
	PUTBUCKET_LOGGING
	PUTBUCKET_NOTIFICATION
	PUTBUCKET_REPLICATION
	PUTBUCKET_TAGGING
	PUTBUCKET_REQUESTPAYMENT
	PUTBUCKET_VERSIONING
	PUTBUCKET_WEBSITE
	HEADBUCKET
	DELETEOBJECT
	GETOBJECT
	GETOBJECT_ACL
	GETOBJECT_TORRENT
	HEADOBJECT
	POSTOBJECT
	POSTOBJECT_RESTORE
	PUTOBJECT
	PUTOBJECT_ACL
	PUTOBJECT_COPY
)

type S3Request struct {
	common.RequestHeaders
	bucket               string
	object               string
	method               string
	s3method             S3METHOD
	contentEncoding      string
	lastModified         time.Time
	versionID            string
	serversideEncryption bool
}

func (s S3Request) String() string {
	return fmt.Sprintf("%s/%s | %s / %v", s.bucket, s.object, s.method, s.s3method)
}
