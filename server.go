package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"strconv"

	"github.com/0x434D53/s3server/common"
	"github.com/0x434D53/s3server/s3backend/inMemory"
)

var port = flag.String("port", "10001", "Server will run on this port")
var hostname = flag.String("host", "test.dev", "Hostname analogous to s3.amazonaws.com")
var basePath = flag.String("basepath", "s3", "Basepath for S3")
var host string

var backend common.S3Backend

func writeError(w http.ResponseWriter, awserr *common.Error) error {
	b, err := xml.Marshal(awserr)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(awserr.StatusCode)

	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	w.Write(b)

	return nil
}

func setCommondResponseHeaders(w http.ResponseWriter, h *common.ResponseHeaders) {
	hd := w.Header()

	if h.ContentLength != "" {
		hd.Set("Content-Length", h.ContentLength)
	}

	if h.ContentType != "" {
		hd.Set("Content-Type", h.ContentType)
	}

	hd.Set("Date", h.Date)
	hd.Set("ETag", h.ETag)

	if h.Server != "" {
		hd.Set("Server", h.Server)
	}

	hd.Set("x-amz-delete-marker", fmt.Sprintf("%v", h.XAmzDeleteMarker))
	hd.Set("x-amz-id-2", h.XAmzId2)
	hd.Set("x-amz-request-id", h.XAmzRequestId)
	hd.Set("x-amz-version-id", h.XAmzVersionId)
}

func logHandlerCall(handler string, rd *S3Request) {
	log.Printf("=====Web===== [%s / %s] %s/%s | %v |", handler, rd.s3method, rd.bucket, rd.object, rd.params)
}

func headBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("headBucketHandler", rd)
	err := backend.HeadBucket(rd.bucket, rd.Authorization)

	log.Printf("HeadBucketHandler: %v", err)

	if err != nil {
		writeError(w, err)
	} else {
		w.WriteHeader(200)
	}
}

func getBucketsHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketsHandler", rd)
}

func putBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("putBucketHandler", rd)
	awserr := backend.PutBucket(rd.bucket, rd.Authorization)

	if awserr != nil {
		writeError(w, awserr)
	} else {
		w.WriteHeader(200)
	}
}

func getBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketHandler", rd)
	lbr, awserr := backend.GetBucketObjects(rd.bucket, rd.Authorization)

	if awserr != nil {
		writeError(w, awserr)
		return
	}

	b, err := xml.Marshal(lbr)

	if err != nil {
		log.Print(err)
		http.Error(w, "Error", 500)
		return
	}

	w.Write(b)
}

func getBucketLocationHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketLocationHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func deleteBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("deleteBucketHandler", rd)
	err := backend.DeleteBucket(rd.bucket, rd.Authorization)

	if err != nil {
	}
	http.Error(w, "Not Implemented", 500)
}

func postObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("postObjectHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func getObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getObjectHandler", rd)
	data, ct, err := backend.GetObject(rd.bucket, rd.object, rd.Authorization)

	if err != nil {
		writeError(w, err)
		return
	}

	rh := &common.ResponseHeaders{
		ContentLength: fmt.Sprintf("%d", len(data)),
		ContentType:   ct,
	}

	setCommondResponseHeaders(w, rh)

	w.Write(data)
}

func putObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("putObjectHandler", rd)
	contents, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, &common.ErrInternalError)
		return
	}

	awserr := backend.PutObject(rd.bucket, rd.object, contents, rd.ContentType, rd.Authorization)

	if awserr != nil {
		writeError(w, awserr)
		return
	}

	w.WriteHeader(200)
}

func copyObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("copyObjectHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func headObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("headObjectHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("deleteObjectHandler", rd)
	err := backend.DeleteObject(rd.bucket, rd.object, "")

	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getBucketObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketObjectVersionHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func getBucketVersioningHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketVersioningHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func getBucketACLHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketACLHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketCORSHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketCORSHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketLifecycleHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketLifecycleHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketPolicyHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketPolicyHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketLoggingHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketLoggingHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketNotificationHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketNotificationHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketReplicationHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketReplicationHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketTagging(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketTagging", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketRequestPaymentHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketRequestPaymentHandler", rd)
	http.Error(w, "Not Implemented", 500)
}
func getBucketWebsiteHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getBucketWebsiteHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func putBucketVersioningHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("putBucketVersioningHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func deleteObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	logHandlerCall("deleteObjectVersionHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func getObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("getObjectVersionHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func headObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	logHandlerCall("headObjectVersionHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func putObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	logHandlerCall("putObjectVersionHandler", rd)
	http.Error(w, "Not Implemented", 500)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	rd, err := getS3RequestData(r)

	if err != nil {
		log.Print(err)
		http.Error(w, "Error parsing the request", 500)

		return
	}

	switch rd.s3method {
	case GETBUCKET:
		getBucketHandler(w, r, rd)
	case GETBUCKET_ACL:
		getBucketACLHandler(w, r, rd)
	case GETBUCKET_CORS:
		getBucketCORSHandler(w, r, rd)
	case GETBUCKET_LIFECYCLE:
		getBucketLifecycleHandler(w, r, rd)
	case GETBUCKET_POLICY:
		getBucketPolicyHandler(w, r, rd)
	case GETBUCKET_LOCATION:
		getBucketLocationHandler(w, r, rd)
	case GETBUCKET_LOGGING:
		getBucketLoggingHandler(w, r, rd)
	case GETBUCKET_NOTIFICATION:
		getBucketNotificationHandler(w, r, rd)
	case GETBUCKET_REPLICATION:
		getBucketReplicationHandler(w, r, rd)
	case GETBUCKET_TAGGING:
		getBucketTagging(w, r, rd)
	case GETBUCKET_OBJECTVERSION:
		getBucketObjectVersionHandler(w, r, rd)
	case GETBUCKET_REQUESTPAYMENT:
		getBucketRequestPaymentHandler(w, r, rd)
	case GETBUCKET_VERSIONING:
		getBucketVersioningHandler(w, r, rd)
	case GETBUCKET_WEBSITE:
		getBucketWebsiteHandler(w, r, rd)
	case HEADBUCKET:
		headBucketHandler(w, r, rd)
	case PUTBUCKET:
		putBucketHandler(w, r, rd)
	case DELETEBUCKET:
		deleteBucketHandler(w, r, rd)
	case PUTBUCKET_ACL:
	case PUTBUCKET_CORS:
	case PUTBUCKET_LIFECYCLE:
	case PUTBUCKET_POLICY:
	case PUTBUCKET_LOGGING:
	case PUTBUCKET_NOTIFICATION:
	case PUTBUCKET_REPLICATION:
	case PUTBUCKET_TAGGING:
	case PUTBUCKET_REQUESTPAYMENT:
	case PUTBUCKET_VERSIONING:
	case DELETEOBJECT:
		deleteObjectHandler(w, r, rd)
	case GETOBJECT:
		getObjectHandler(w, r, rd)
	case GETOBJECT_ACL:
	case GETOBJECT_TORRENT:
	case HEADOBJECT:
		headObjectHandler(w, r, rd)
	case POSTOBJECT:
		postObjectHandler(w, r, rd)
	case POSTOBJECT_RESTORE:
	case PUTOBJECT:
		putObjectHandler(w, r, rd)
	case PUTOBJECT_ACL:
	case PUTOBJECT_COPY:
	default:
		http.Error(w, "Unkown Error", 200)
	}
}

func getS3RequestData(r *http.Request) (*S3Request, *common.Error) {
	s3r := S3Request{}

	pathMethod := false

	if r.Host == host {
		pathMethod = true
	}

	if pathMethod {
		s := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(s) == 0 {
			return nil, &common.ErrInvalidBucketName
		} else if len(s) == 1 {
			s3r.bucket = s[0]
		} else if len(s) == 2 {
			s3r.bucket = s[0]
			s3r.object = s[1]
		} else {
			return nil, &common.ErrInvalidArgument
		}
	} else {
		//TODO
	}

	s3r.params = r.URL.Query()

	if s3r.object == "" {
		switch r.Method {
		case "POST":
			return nil, &common.ErrMethodNotAllowed
		case "PUT":
			s3r.s3method = PUTBUCKET
		case "DELETE":
			s3r.s3method = DELETEBUCKET
		case "HEAD":
			s3r.s3method = HEADBUCKET
		case "GET":
			if s3r.HasParam("location") {
				s3r.s3method = GETBUCKET_LOCATION
			} else if s3r.HasParam("acl") {
				s3r.s3method = GETBUCKET_ACL
			} else if s3r.HasParam("cors") {
				s3r.s3method = GETBUCKET_CORS
			} else if s3r.HasParam("lifecycle") {
				s3r.s3method = GETBUCKET_LIFECYCLE
			} else if s3r.HasParam("policy") {
				s3r.s3method = GETBUCKET_POLICY
			} else if s3r.HasParam("logging") {
				s3r.s3method = GETBUCKET_LOGGING
			} else if s3r.HasParam("notification") {
				s3r.s3method = GETBUCKET_NOTIFICATION
			} else if s3r.HasParam("replication") {
				s3r.s3method = GETBUCKET_REPLICATION
			} else if s3r.HasParam("tagging") {
				s3r.s3method = GETBUCKET_TAGGING
			} else if s3r.HasParam("versions") {
				s3r.s3method = GETBUCKET_VERSIONING
			} else if s3r.HasParam("requestPayment") {
				s3r.s3method = GETBUCKET_REQUESTPAYMENT
			} else if s3r.HasParam("versioning") {
				s3r.s3method = GETBUCKET_VERSIONING
			} else if s3r.HasParam("website") {
				s3r.s3method = GETBUCKET_WEBSITE
			} else {
				s3r.s3method = GETBUCKET
			}
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

	cl := r.Header.Get("Content-Length")

	if cl != "" {
		cli, err := strconv.ParseUint(cl, 10, 64)
		if err != nil {
			return nil, &common.ErrInvalidRequest
		}
		s3r.ContentLength = cli
	}
	s3r.ContentType = r.Header.Get("Content-Type")
	s3r.contentEncoding = r.Header.Get("Content-Encoding")
	s3r.ContentMD5 = r.Header.Get("Content-MD5")

	return &s3r, nil
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=====RESET=====")
	backend.Reset()
}

func main() {
	flag.Parse()
	host = *hostname + ":" + *port
	//	host = *hostname + ":" + *port
	backend = inMemory.NewS3Backend()

	fmt.Printf("Launching S3Server on port %v\n", *port)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/_internal/reset", resetHandler)

	log.Fatal(http.ListenAndServe(":"+string(*port), nil))
}
