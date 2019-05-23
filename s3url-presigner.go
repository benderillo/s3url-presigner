package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/benderillo/s3url-presigner/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Verbose         bool   `short:"v" long:"verbose" description:"Print detailed information"`
	AwsRegion       string `short:"r" long:"aws-region" description:"AWS region" env:"AWS_REGION" required:"true"`
	AwsAccessID     string `short:"i" long:"aws-access-id" description:"AWS access ID" env:"AWS_ACCESS_KEY_ID" required:"true"`
	AwsSecretKey    string `short:"s" long:"aws-secret-key" description:"AWS secret key" env:"AWS_SECRET_ACCESS_KEY" required:"true"`
	AwsSessionToken string `short:"t" long:"aws-session-token" description:"AWS session token" env:"AWS_SESSION_TOKEN"`
	Bucket          string `short:"b" long:"bucket" description:"S3 bucket" env:"S3_BUCKET" required:"true"`
	Path            string `short:"p" long:"path" description:"S3 path" env:"S3_PATH" required:"true"`
	Method          string `short:"m" long:"method" choice:"get" choice:"put" description:"HTTP method that needs to be presigned" default:"get"`
	Expiry          int64  `short:"e" long:"expiry" description:"Expiration time for the url in seconds" default:"7200"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	setVerboseLevel(opts.Verbose)

	// Create AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(opts.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			opts.AwsAccessID,
			opts.AwsSecretKey,
			opts.AwsSessionToken,
		),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create AWS session")
	}

	s3 := storage.NewStorage(awsSession)

	urlStr := "s3://" + opts.Bucket + opts.Path
	expiry := time.Second * time.Duration(opts.Expiry)

	s3Url, expTime, err := s3.GetPresignedURL(opts.Method, urlStr, expiry)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to generate presigned url")
	}

	output, _ := json.MarshalIndent(struct {
		URL    string
		Expiry time.Time
	}{URL: *s3Url, Expiry: expTime}, "", "\t")

	fmt.Print(string(output))
	os.Exit(0)
}

func setVerboseLevel(verbose bool) {
	var logWriter io.Writer
	if verbose {
		// Do not print logs as JSON
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		logWriter = os.Stdout
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = zerolog.New(logWriter).With().Timestamp().Logger()
}
