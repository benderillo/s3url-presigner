# s3url-presigner
[![Go Report Card](https://goreportcard.com/badge/github.com/benderillo/s3url-presigner)](https://goreportcard.com/report/github.com/benderillo/s3url-presigner)
[![GolangCI](https://golangci.com/badges/github.com/benderillo/s3url-presigner.svg)](https://golangci.com/r/github.com/benderillo/s3url-presigner)
[![Release](https://img.shields.io/github/release/benderillo/s3url-presigner.svg)](https://github.com/benderillo/s3url-presigner/releases/latest)

S3 URL presigner generates pre-signed S3 urls for PUT and GET requests

### To install to your $GOPATH/bin
 `go get github.com/benderillo/s3url-presigner/cmd/s3url-presigner`

### Usage

```
Usage:
  aws-s3-presign-url [OPTIONS]

Application Options:
  -r, --aws-region=        AWS region [$AWS_REGION]
  -i, --aws-access-id=     AWS access ID [$AWS_ACCESS_KEY_ID]
  -s, --aws-secret-key=    AWS secret key [$AWS_SECRET_ACCESS_KEY]
  -t, --aws-session-token= AWS session token [$AWS_SESSION_TOKEN]
  -b, --bucket=            S3 bucket [$S3_BUCKET]
  -p, --path=              S3 path [$S3_PATH]
  -m, --method=[get|put]   HTTP method that needs to be presigned (default: get)
  -e, --expiry=            Expiration time for the url in seconds (default: 7200)

Help Options:
  -h, --help               Show this help message
  ```

### Example input

```
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=myaccesskeyid
export AWS_SECRET_ACCESS_KEY=mysecretaccesskey

s3url-presigner --bucket my-example-bucket --path /test/path/file.txt --method put --expiry 3600
```

### Example output
```
{
	"URL": "https://my-example-bucket.s3.us-east-1.amazonaws.com/test/path/file.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=XXXXXXXXXXXXX%2F20190523%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20190523T032122Z\u0026X-Amz-Expires=3600\u0026X-Amz-Security-Token=FjopijpoirjpeoirgjsofdighsdfoighdiohgXXXXXXXXXXXXxxxxxxxxxxx\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=iddqdiddqdgodshowmethemoneygodiddqdiddqd",
	"Expiry": "2019-05-23T04:21:22.919215Z"
}
```

### How to include it in your project
```
import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/benderillo/s3url-presigner/pkg/storage"

	"golang.org/x/xerrors"
)

func main() {
	awsSession := session.Must(NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	s3 := storage.NewStorage(awsSession)

	presignedUrl, expTime, err := s3.GetPresignedURL("put", "s3://bucket/path", time.Hour)

}
```
