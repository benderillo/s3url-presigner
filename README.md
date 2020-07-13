# s3url-presigner
[![Go Report Card](https://goreportcard.com/badge/github.com/benderillo/s3url-presigner)](https://goreportcard.com/report/github.com/benderillo/s3url-presigner)
[![GolangCI](https://golangci.com/badges/github.com/benderillo/s3url-presigner.svg)](https://golangci.com/r/github.com/benderillo/s3url-presigner)
[![Release](https://img.shields.io/github/release/benderillo/s3url-presigner.svg)](https://github.com/benderillo/s3url-presigner/releases/latest)
[![GoDoc](https://godoc.org/github.com/benderillo/s3url-presigner?status.svg)](https://godoc.org/github.com/benderillo/s3url-presigner/pkg/storage)

S3 URL presigner generates pre-signed S3 urls for PUT and GET requests

### To install to your $GOPATH/bin
 `go get github.com/benderillo/s3url-presigner/cmd/s3url-presigner`

### Usage

```
Usage:
  main [OPTIONS]

Application Options:
  -u, --url=             S3 URL (s3://bucket/path) [$S3_URL]
  -m, --method=[get|put] HTTP method that needs to be presigned (default: get)
  -e, --expiry=          Expiration time for the url in seconds (default: 7200)

Help Options:
  -h, --help             Show this help message
```

### Example input

```
s3url-presigner --url s3://my-example-bucket/test/path/file.txt --method put --expiry 3600
```

### Example output
```
https://my-example-bucket.s3.us-east-1.amazonaws.com/test/path/file.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=XXXXXXXXXXXXX%2F20190523%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20190523T032122Z\u0026X-Amz-Expires=3600\u0026X-Amz-Security-Token=FjopijpoirjpeoirgjsofdighsdfoighdiohgXXXXXXXXXXXXxxxxxxxxxxx\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=iddqdiddqdgodshowmethemoneygodiddqdiddqd
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
