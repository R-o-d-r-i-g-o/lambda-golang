package internal

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Handler interface {
	Exec(ctx context.Context) (string, error)
}

type handler struct {
	sqsClient *sqs.SQS
	queueURL  string
}

func NewHandler(sqsClient *sqs.SQS, queueURL string) Handler {
	return &handler{sqsClient, queueURL}
}

// Exec plays lambda action.
func (h *handler) Exec(ctx context.Context) (string, error) {
	start := time.Now()

	// TODO: add consumer & processors layer.

	// take SQS messagens from fifo queue.
	result, err := h.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(h.queueURL),
		MaxNumberOfMessages: aws.Int64(MAX_NUM_OF_MESSAGES),
		WaitTimeSeconds:     aws.Int64(WAIT_TIME_SECONDS),
	})
	if err != nil {
		log.Errorf("failed receive msgs. Error: %s", err.Error())
		return "", err
	}

	log.Info("get messages correclty.")

	var (
		lastSenderId = ""
		hasChanged   = false
	)

	for _, msg := range result.Messages {
		res, err := unmarshalResponse[Message](*msg.Body)
		if err != nil {
			log.Warnf("failed unmarshal msg. Error: %s", err.Error())
			continue
		}
		if lastSenderId == "" {
			lastSenderId = res.Sender.ID
			continue
		}

		hasChanged = lastSenderId != res.Sender.ID
	}

	log.Info("finished checks correctly.")

	// Note: don't delete messages this step.
	res, err := formatResponseToSNS(len(result.Messages), start, hasChanged)
	return string(res), err
}

// formatResponseToSNS to be send as a report to SNS topics.
func formatResponseToSNS(msgLen int, execTime time.Time, hasChanged bool) ([]byte, error) {
	status := struct {
		MessagesProcessed int           `json:"amount_of_messages_processed"`
		Duration          time.Duration `json:"duration"`
		HasChanged        bool          `json:"user_has_changed"`
	}{
		MessagesProcessed: msgLen,
		Duration:          time.Since(execTime),
		HasChanged:        hasChanged,
	}

	return marshalElement(status)
}
