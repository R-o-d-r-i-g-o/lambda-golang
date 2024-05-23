package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"lambda-golang/internal"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

/**
 * The purpose of this lambda is check for same sender, then notify it in SNS service.
 * This integration is needed for a POC into "Ordery App", by usin' sqs message protocol.
 */

var (
	sqsClient *sqs.SQS
	queueURL  string
)

func init() {
	// format log messages.
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Sign in into aws sdk service.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqsClient = sqs.New(sess)
	queueURL = os.Getenv("SQS_QUEUE_URL")

	log.Info("finished sign in aws sdk.")
}

func main() {
	handle := internal.NewHandler(sqsClient, queueURL)
	lambda.Start(handle.Exec)
}
