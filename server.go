package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/0x434D53/s3server/common"
	"github.com/0x434D53/s3server/s3backend/inMemory"
	"strconv"
)

var port = flag.String("port", "10001", "Server will run on this port")
var hostname = flag.String("host", "localhost", "Hostname analogous to w3.amazonaws.com")
var basePath = flag.String("basepath", "s3", "Basepath for S3")
var host string

var backend common.S3Backend

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

func WriteAndLogError(err error, w http.ResponseWriter, r *http.Request, rd *S3Request) {
	log.Print(err)
	if err == common.ErrNotFound {
		http.Error(w, "Not Found", 404)
	} else if err == common.ErrAlreadyExists {
		http.Error(w, "Already exists", 409)
	} else {
		http.Error(w, "Unknown error", 500)
	}
}

func headBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	err := backend.HeadBucket(rd.bucket, rd.Authorization)

	log.Printf("HeadBucketHandler: %v", err)

	if err != nil {
		WriteAndLogError(err, w, r, rd)
	} else {
		w.WriteHeader(200)
	}
}

func getBucketsHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

}

func putBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	err := backend.PutBucket(rd.bucket, rd.Authorization)
	log.Printf("PutBucketHandler: %v", err)

	if err != nil {
		WriteAndLogError(err, w, r, rd)
	} else {
		w.WriteHeader(200)
	}
}

func getBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	lbr, err := backend.GetBucketObjects(rd.bucket, rd.Authorization)

	if err == common.ErrNotFound {
		http.Error(w, "Not Found", 404)
		return
	} else if err != nil {
		log.Print(err)
		http.Error(w, "Error", 500)
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
	http.Error(w, "Not Implemented", 500)
}

func deleteBucketHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	http.Error(w, "Not Implemented", 500)
}

func postObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func getObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	data, ct, err := backend.GetObject(rd.bucket, rd.object, rd.Authorization)

	if err != nil {
		if err == common.ErrNotFound {
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, "InternalServeError", 500)
			return
		}
	}

	rh := &common.ResponseHeaders{
		ContentLength: fmt.Sprintf("%d", len(data)),
		ContentType: ct,
	}

	setCommondResponseHeaders(w, rh)

	w.Write(data)
}

func putObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	contents, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Unkown Error", 500)
		return
	}

	err = backend.PutObject(rd.bucket, rd.object, contents, rd.ContentType, rd.Authorization)

	if err != nil {
		http.Error(w, "Not Implemented", 500)
	}

	w.WriteHeader(200)
}

func copyObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	http.Error(w, "Not Implemented", 500)
}

func headObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	http.Error(w, "Not Implemented", 500)
}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {
	err := backend.DeleteObject(rd.bucket, rd.object, "")

	if err != nil {
		http.Error(w, "Not Implemented", 500)

	}

	w.WriteHeader(http.StatusOK)
}

func getBucketObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func getBucketVersioningHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func putBucketVersioningHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func deleteObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func getObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func headObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func putObjectVersionHandler(w http.ResponseWriter, r *http.Request, rd *S3Request) {

	http.Error(w, "Not Implemented", 500)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	rd, err := getS3RequestData(r)

	if err != nil {
		log.Print(err)
		http.Error(w, "Error parsing the request", 500)

		return
	}

	log.Print(rd.s3method)

	switch rd.s3method {
	case GETBUCKET:
		getBucketHandler(w, r, rd)
	case GETBUCKET_ACL:
	case GETBUCKET_CORS:
	case GETBUCKET_LIFECYCLE:
	case GETBUCKET_POLICY:
	case GETBUCKET_LOCATION:
		getBucketLocationHandler(w, r, rd)
	case GETBUCKET_LOGGING:
	case GETBUCKET_NOTIFICATION:
	case GETBUCKET_REPLICATION:
	case GETBUCKET_TAGGING:
	case GETBUCKET_OBJECTVERSION:
	case GETBUCKET_REQUESTPAYMENT:
	case GETBUCKET_VERSIONING:
		getBucketVersioningHandler(w, r, rd)
	case GETBUCKET_WEBSITE:
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

	fmt.Printf("%v\n", rd)
}

func getS3RequestData(r *http.Request) (*S3Request, error) {
	s3r := S3Request{}

	pathMethod := false

	if r.Host == host {
		pathMethod = true
	}

	log.Print(r.URL.Path)

	if pathMethod {
		s3r.bucket, s3r.object = path.Split(r.URL.Path)
		s3r.bucket = strings.Trim(s3r.bucket, "/")
	} else {
		//TODO
	}

	log.Print(r.Method)

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

	cl :=r.Header.Get("Content-Length")

	if cl != "" {
		cli, err := strconv.ParseUint(cl, 10, 64)
		if err != nil {
			return nil, err
		}
		s3r.ContentLength = cli
	}
	s3r.ContentType = r.Header.Get("Content-Type")
	s3r.contentEncoding = r.Header.Get("Content-Encoding")
	s3r.ContentMD5 = r.Header.Get("Content-MD5")

	return &s3r, nil
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("--- RESET ---")
	backend.Reset()
}

func main() {
	flag.Parse()
	host = *hostname + ":" + *port
	backend = inMemory.NewS3Backend()

	fmt.Printf("Launching S3Server on port %v\n", *port)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/_internal/reset", resetHandler)

	log.Fatal(http.ListenAndServe(host, nil))
}
