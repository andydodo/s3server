package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/0x434D53/s3server/s3backend"
	"github.com/0x434D53/s3server/s3backend/common"
)

const (
	HOST     = "localhost:10001"
	BASEPATH = "s3"
)

type S3METHOD int

type ListResp struct {
	Name       string
	Prefix     string
	Delimiter  string
	Marker     string
	NextMarker string
	MaxKeys    int

	// IsTruncated is true if the results have been truncated because
	// there are more keys and prefixes than can fit in MaxKeys.
	// N.B. this is the opposite sense to that documented (incorrectly) in
	// http://goo.gl/YjQTc
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
	HEADBUCKET
	DELETEOBJECT
	GETOBJECT
	HEADOBJECT
	POSTOBJECT
	PUTOBJECT
)

type S3Request struct {
	bucket               string
	object               string
	method               string
	s3method             S3METHOD
	authorization        string
	contentMD5           []byte
	date                 time.Time
	contentLength        string
	contentEncoding      string
	lastModified         time.Time
	versionId            string
	deleteMarker         bool
	serversideEncryption bool
}

var backend s3backend.S3Backend

func getBucketsHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func putBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
}

func getBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func getBucketLocationHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func deleteBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func postObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func getObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func putObjectHandler(w http.ResponseWriter, r *http.Request) {

}

func copyObjectHandler(w http.ResponseWriter, r *http.Request) {

}

func headObjectHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request) {

}

func getBucketObjectVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func getBucketVersioningHandler(w http.ResponseWriter, r *http.Request) {

}

func putBucketVersioningHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteObjectVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func getObjectVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func headObjectVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func putObjectVersionHandler(w http.ResponseWriter, r *http.Request) {

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	rd, err := getS3RequestData(r)

	if err != nil {
		log.Print(err)
		http.Error(w, "Error parsing the request", 500)

		return
	}

	switch rd.s3method {
	case PUTBUCKET:
		putBucketHandler(w, r, rd)
	case GETBUCKET
	}

	fmt.Printf("%v\n", rd)
}

func getS3RequestData(r *http.Request) (*S3Request, error) {
	s3r := S3Request{}

	pathMethod := false

	if r.Host == HOST {
		pathMethod = true
		fmt.Println("PATH")
	}

	if pathMethod {
		s3r.bucket, s3r.object = path.Split(r.URL.Path)
		s3r.bucket = strings.Trim(s3r.bucket, "/")
	} else {
		//TODO
	}

	if s3r.object == "" {
		switch r.Method {
		case "POST":
			return nil, fmt.Errorf("No POST on Buckets")
		case "PUT":
			s3r.s3method = PUTBUCKET
		case "DELETE":
			s3r.s3method = DELETEBUCKET
		case "HEAD":
			s3r.s3method = HEADBUCKET
		case "GET":
			s3r.s3method = GETBUCKET
		}
	} else {
		switch r.Method {
		case "POST":
			s3r.s3method = POSTOBJECT
		case "PUT":
			s3r.s3method = PUTOBJECT
		case "DELETE":
			s3r.s3method = DELETEOBJECT
		case "HEAD":
			s3r.s3method = HEADOBJECT
		case "GET":
			s3r.s3method = GETOBJECT
		}
	}

	s3r.contentLength = r.Header.Get("Content-Length")
	s3r.contentEncoding = r.Header.Get("Content-Encoding")
	s3r.contentMD5 = []byte(r.Header.Get("Content-MD5"))

	return &s3r, nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Println(r.Host)
	fmt.Println(r.Header)
	fmt.Println(r.Method)

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Print(err)
		return
	}

	fmt.Println(string(b))
	w.WriteHeader(200)
}

func main() {
	backend = s3backend.NewS3InMemory()

	http.HandleFunc("/", mainHandler)

	log.Fatal(http.ListenAndServe(HOST, nil))
}