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

	"strconv"

	"github.com/0x434D53/s3server/common"
	"github.com/0x434D53/s3server/s3backend/inMemory"
)

var port = flag.String("port", "10001", "Server will run on this port")
var hostname = flag.String("host", "localhost", "Hostname analogous to w3.amazonaws.com")
var basePath = flag.String("basepath", "s3", "Basepath for S3")
var host string

var backend common.S3Backend

var (
	ErrAccessDenied                            = Error{http.StatusForbidden, "AccessDenied", "Access Denied", "", "", ""}
	ErrAccountProblem                          = Error{http.StatusForbidden, "AccountProblem", "There is a problem with your AWS account that prevents the operation from completing successfully.", "", "", ""}
	ErrAmbigiousGrantByEmailAddress            = Error{http.StatusBadRequest, "AmbiguousGrantByEmailAddress", "The email address you provided is associated with more than one account.", "", "", ""}
	ErrBadDigest                               = Error{http.StatusBadRequest, "BadDigest", "The Content-MD5 you specified did not match what we received.", "", "", ""}
	ErrBucketAlreadyExists                     = Error{http.StatusConflict, "BucketAlreadyExists", "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.", "", "", ""}
	ErrBucketAlreadyOwnedByYou                 = Error{http.StatusConflict, "BucketAlreadyOwnedByYou", "Your previous request to create the named bucket succeeded and you already own it. You get this error in all AWS regions except US Standard, us-east-1. In us-east-1 region, you will get 200 OK, but it is no-op (if bucket exists it Amazon S3 will not do anything).", "", "", ""}
	ErrBucketNotEmpty                          = Error{http.StatusConflict, "BucketNotEmpty", "The bucket you tried to delete is not empty.", "", "", ""}
	ErrCredentialsNotSupported                 = Error{http.StatusBadRequest, "CredentialsNotSupported", "This request does not support credentials.	", "", "", ""}
	ErrCrossLocationLoggingProhibited          = Error{http.StatusForbidden, "CrossLocationLoggingProhibited", "Cross-location logging not allowed. Buckets in one geographic location cannot log information to a bucket in another location.	", "", "", ""}
	ErrEntityTooSmall                          = Error{http.StatusBadRequest, "EntityTooSmall", "Your proposed upload is smaller than the minimum allowed object size.	", "", "", ""}
	ErrEntityTooLarge                          = Error{http.StatusBadRequest, "EntityTooLarge", "Your proposed upload exceeds the maximum allowed object size.", "", "", ""}
	ErrExpiredToken                            = Error{http.StatusBadRequest, "ExpiredToken", "The provided token has expired.", "", "", ""}
	ErrIllegalVersioningConfigurationException = Error{http.StatusBadRequest, "IllegalVersioningConfigurationException", "Indicates that the versioning configuration specified in the request is invalid.", "", "", ""}
	ErrIncompleteBody                          = Error{http.StatusBadRequest, "IncompleteBody", "You did not provide the number of bytes specified by the Content-Length HTTP header	", "", "", ""}
	ErrIncorrectNumberOfFilesInPostRequest     = Error{http.StatusBadRequest, "IncorrectNumberOfFilesInPostRequest", "POST requires exactly one file upload per request.", "", "", ""}
	ErrInlineDataTooLarge                      = Error{http.StatusBadRequest, "InlineDataTooLarge", "Inline data exceeds the maximum allowed size.", "", "", ""}
	ErrInternalError                           = Error{http.StatusInternalServerError, "InternalError", "We encountered an internal error. Please try again.", "", "", ""}
	ErrInvalidAccessKeyId                      = Error{http.StatusForbidden, "InvalidAccessKeyId", "The AWS access key Id you provided does not exist in our records.", "", "", ""}
	ErrInvalidAddressingHeader                 = Error{0, "InvalidAddressingHeader", "You must specify the Anonymous role.", "", "", ""} // Unclear HTTP Status Code???
	ErrInvalidArgument                         = Error{http.StatusBadRequest, "InvalidArgument", "Invalid Argument", "", "", ""}
	ErrInvalidBucketName                       = Error{http.StatusBadRequest, "InvalidBucketName", "The specified bucket is not valid.", "", "", ""}
	ErrInvalidBucketState                      = Error{http.StatusConflict, "InvalidBucketState", "The request is not valid with the current state of the bucket.", "", "", ""}
	ErrInvalidDigest                           = Error{http.StatusBadRequest, "InvalidDigest", "The Content-MD5 you specified is not valid.", "", "", ""}
	ErrInvalidEncryptionAlgorithmError         = Error{http.StatusBadRequest, "InvalidEncryptionAlgorithmError", "The encryption request you specified is not valid. The valid value is AES256", "", "", ""}
	ErrInvalidLocationConstraint               = Error{http.StatusBadRequest, "InvalidLocationConstraint", "", "", "", ""}
	ErrInvalidObjectState                      = Error{http.StatusForbidden, "InvalidObjectState", "The operation is not valid for the current state of the object.	", "", "", ""}
	ErrInvalidPart                             = Error{http.StatusBadRequest, "InvalidPart", "One or more of the specified parts could not be found. The part might not have been uploaded, or the specified entity tag might not have matched the part's entity tag.	", "", "", ""}
	ErrInvalidPartOrder                        = Error{http.StatusBadRequest, "InvalidPartOrder", "The list of parts was not in ascending order.Parts list must specified in order by part number.	", "", "", ""}
	ErrInvalidPayer                            = Error{http.StatusForbidden, "InvalidPayer", "All access to this object has been disabled.", "", "", ""}
	ErrInvalidPolicyDocument                   = Error{http.StatusBadRequest, "InvalidPolicyDocument", "The content of the form does not meet the conditions specified in the policy document.	", "", "", ""}
	ErrInvalidRange                            = Error{http.StatusRequestedRangeNotSatisfiable, "InvalidRange", "The requested range cannot be satisfied.	", "", "", ""}
	ErrInvalidSecurity                         = Error{http.StatusBadRequest, "InvalidSecurity", "The provided security credentials are not valid.	", "", "", ""}
	ErrInvalidSOAPRequest                      = Error{http.StatusBadRequest, "InvalidSOAPRequest", "The SOAP request body is invalid.	", "", "", ""}
	ErrInvalidRequest                          = Error{http.StatusBadRequest, "InvalidRequest", "SOAP requests must be made over an HTTPS connection", "", "", ""}
	ErrInvalidStorageClass                     = Error{http.StatusBadRequest, "InvalidStorageClass", "The storage class you specified is not valid.	", "", "", ""}
	ErrInvalidTargetBucketForLogging           = Error{http.StatusBadRequest, "InvalidTargetBucketForLogging", "The target bucket for logging does not exist, is not owned by you, or does not have the appropriate grants for the log-delivery group.	", "", "", ""}
	ErrInvalidToken                            = Error{http.StatusBadRequest, "InvalidToken", "The provided token is malformed or otherwise invalid.	", "", "", ""}
	ErrInvalidURI                              = Error{http.StatusBadRequest, "InvalidURI", "Couldn't parse the specified URI.	", "", "", ""}
	ErrKeyTooLong                              = Error{http.StatusBadRequest, "KeyTooLong", "Your key is too long.	", "", "", ""}
	ErrMalformedACLError                       = Error{http.StatusBadRequest, "MalformedACLError", "The XML you provided was not well-formed or did not validate against our published schema.	", "", "", ""}
	ErrMalformedPOSTRequest                    = Error{http.StatusBadRequest, "MalformedPOSTRequest", "The body of your POST request is not well-formed multipart/form-data.	", "", "", ""}
	ErrMalformedXML                            = Error{http.StatusBadRequest, "MalformedXML", "The XML you provided was not well-formed or did not validate against our published schema.", "", "", ""}
	ErrMaxMessageLengthExceeded                = Error{http.StatusBadRequest, "MaxMessageLengthExceeded", "Your request was too big.	", "", "", ""}
	ErrMaxPostPreDataLengthExceededError       = Error{http.StatusBadRequest, "MaxPostPreDataLengthExceededError", "Your POST request fields preceding the upload file were too large.	", "", "", ""}
	ErrMetadataTooLarge                        = Error{http.StatusBadRequest, "MetadataTooLarge", "Your metadata headers exceed the maximum allowed metadata size.	", "", "", ""}
	ErrMethodNotAllowed                        = Error{http.StatusMethodNotAllowed, "MethodNotAllowed", "he specified method is not allowed against this resource.", "", "", ""}
	ErrMissingAttachment                       = Error{0, "MissingAttachment", "A SOAP attachment was expected, but none were found.", "", "", ""} //???
	ErrMissingContentLength                    = Error{http.StatusLengthRequired, "MissingContentLength", "You must provide the Content-Length HTTP header.", "", "", ""}
	ErrMissingRequestBodyError                 = Error{http.StatusBadRequest, "MissingRequestBodyError", "Request body is empty.", "", "", ""}
	ErrMissingSecurityElement                  = Error{http.StatusBadRequest, "MissingSecurityElement", "The SOAP 1.1 request is missing a security element", "", "", ""}
	ErrMissingSecurityHeader                   = Error{http.StatusBadRequest, "MissingSecurityHeader", "Your request is missing a required header.	", "", "", ""}
	ErrNoLoggingStatusForKey                   = Error{http.StatusBadRequest, "NoLoggingStatusForKey", "There is no such thing as a logging status subresource for a key.	", "", "", ""}
	ErrNoSuchBucket                            = Error{http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.	", "", "", ""}
	ErrNoSuchKey                               = Error{http.StatusNotFound, "NoSuchKey", "The specified key does not exist.", "", "", ""}
	ErrNoSuchLifecycleConfiguration            = Error{http.StatusNotFound, "NoSuchLifecycleConfiguration", "The lifecycle configuration does not exist.	", "", "", ""}
	ErrNoSuchUpload                            = Error{http.StatusNotFound, "NoSuchUpload", "The specified multipart upload does not exist. The upload ID might be invalid, or the multipart upload might have been aborted or completed.	", "", "", ""}
	ErrNoSuchVersion                           = Error{http.StatusNotFound, "NoSuchVersion", "Indicates that the version ID specified in the request does not match an existing version", "", "", ""}
	ErrNotImplemented                          = Error{http.StatusNotImplemented, "NotImplemented", "A header you provided implies functionality that is not implemented.", "", "", ""}
	ErrNotSignedUp                             = Error{http.StatusForbidden, "NotSignedUp", "Your account is not signed up for the Amazon S3 service. You must sign up before you can use Amazon S3. ", "", "", ""}
	ErrNoSuchBucketPolicy                      = Error{http.StatusForbidden, "NoSuchBucketPolicy", "The specified bucket does not have a bucket policy.", "", "", ""}
	ErrOperationAborted                        = Error{http.StatusConflict, "OperationAborted", "A conflicting conditional operation is currently in progress against this resource. Try again.", "", "", ""}
	ErrPermanentRedirect                       = Error{http.StatusMovedPermanently, "PermanentRedirect", "The bucket you are attempting to access must be addressed using the specified endpoint. Send all future requests to this endpoint.", "", "", ""}
	ErrPreconditionFailed                      = Error{http.StatusPreconditionFailed, "PreconditionFailed", "At least one of the preconditions you specified did not hold.", "", "", ""}
	ErrRedirect                                = Error{307, "Redirect", "Temporary redirect.	", "", "", ""}
	ErrRestoreAlreadyInProgress                = Error{http.StatusConflict, "RestoreAlreadyInProgress", "Object restore is already in progress.", "", "", ""}
	ErrRequestIsNotMultiPartContent            = Error{http.StatusBadRequest, "RequestIsNotMultiPartContent", "Bucket POST must be of the enclosure-type multipart/form-data.	", "", "", ""}
	ErrRequestTimeout                          = Error{http.StatusBadRequest, "RequestTimeout", "Your socket connection to the server was not read from or written to within the timeout period.", "", "", ""}
	ErrRequestTimeTooSkewed                    = Error{http.StatusForbidden, "RequestTimeTooSkewed", "The difference between the request time and the server's time is too large.", "", "", ""}
	ErrRequestTorrentOfBucketError             = Error{http.StatusBadRequest, "RequestTorrentOfBucketError", "Requesting the torrent file of a bucket is not permitted.	", "", "", ""}
	ErrSignatureDoesNotMatch                   = Error{http.StatusForbidden, "SignatureDoesNotMatch", "The request signature we calculated does not match the signature you provided. Check your AWS secret access key and signing method", "", "", ""}
	ErrServiceUnavailable                      = Error{http.StatusServiceUnavailable, "ServiceUnavailable", "Reduce your request rate.	", "", "", ""}
	ErrSlowDown                                = Error{http.StatusServiceUnavailable, "SlowDown", "Reduce your request rate.", "", "", ""}
	ErrTemporaryRedirect                       = Error{307, "TemporaryRedirect", "You are being redirected to the bucket while DNS updates.	", "", "", ""}
	ErrTokenRefreshRequired                    = Error{http.StatusBadRequest, "TokenRefreshRequired", "The provided token must be refreshed.", "", "", ""}
	ErrTooManyBuckets                          = Error{http.StatusBadRequest, "TooManyBuckets", "You have attempted to create more buckets than allowed.", "", "", ""}
	ErrUnexpectedContent                       = Error{http.StatusBadRequest, "UnexpectedContent", "This request does not support content.", "", "", ""}
	ErrUnresolvableGrantByEmailAddress         = Error{http.StatusBadRequest, "UnresolvableGrantByEmailAddress", "The email address you provided does not match any account on record.", "", "", ""}
	ErrUserKeyMustBeSpecified                  = Error{http.StatusBadRequest, "UserKeyMustBeSpecified", "The bucket POST must contain the specified field name. If it is specified, check the order of the fields.	", "", "", ""}
)

func writeError(w http.ResponseWriter, err Error) {

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
	err := backend.DeleteBucket(rd.bucket, rd.Authorization)

	if err != nil {
	}
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
		ContentType:   ct,
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

	cl := r.Header.Get("Content-Length")

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
