package main

import (
	"fmt"
	"os"
	"time"

	"github.com/benderillo/s3url-presigner/pkg/storage"

	"github.com/aws/aws-sdk-go/aws/session"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Bucket string `short:"b" long:"bucket" description:"S3 bucket" env:"S3_BUCKET" required:"true"`
	Path   string `short:"p" long:"path" description:"S3 path" env:"S3_PATH" required:"true"`
	Method string `short:"m" long:"method" choice:"get" choice:"put" description:"HTTP method that needs to be presigned" default:"get"`
	Expiry int64  `short:"e" long:"expiry" description:"Expiration time for the url in seconds" default:"7200"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// Create AWS session
	awsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create AWS session")
	}

	s3 := storage.NewStorage(awsSession)

	urlStr := "s3://" + opts.Bucket + opts.Path
	expiry := time.Second * time.Duration(opts.Expiry)

	s3Url, _, err := s3.GetPresignedURL(opts.Method, urlStr, expiry)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to generate presigned url")
	}

	fmt.Println(*s3Url)
	os.Exit(0)
}
