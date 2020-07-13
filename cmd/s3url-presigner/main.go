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
	Url    string `short:"u" long:"url" description:"S3 URL (s3://bucket/path)" env:"S3_URL" required:"true"`
	Method string `short:"m" long:"method" choice:"get" choice:"put" description:"HTTP method that needs to be presigned" default:"get"`
	Expiry int64  `short:"e" long:"expiry" description:"Expiration time for the url in seconds" default:"7200"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
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

	expiry := time.Second * time.Duration(opts.Expiry)

	s3Url, _, err := s3.GetPresignedURL(opts.Method, opts.Url, expiry)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to generate presigned url")
	}

	fmt.Println(*s3Url)
	os.Exit(0)
}
