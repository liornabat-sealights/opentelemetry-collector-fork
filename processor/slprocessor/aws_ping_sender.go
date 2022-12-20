package slprocessor

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

type AwsPingSender struct {
	awsConnector *AwsConnector
	logger       *zap.Logger
	queueUri     string
	sqsClient    *sqs.SQS
}

func NewAwsPingSender() PingSenderInterface {
	return &AwsPingSender{}
}

func (a *AwsPingSender) Init(logger *zap.Logger, queueUri string) {
	a.logger = logger
	a.queueUri = queueUri

	if a.queueUri == "" {
		a.logger.Error(fmt.Sprintf("Queue URI cannot be empty"))
		return
	}

	a.logger.Info(fmt.Sprintf("Queue URI set to: %s", a.queueUri))

	a.awsConnector = NewAwsConnector()
	err := a.awsConnector.init(a.logger)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Cannot initialize AWS Connector %v", err))
	}

	a.sqsClient = sqs.New(a.awsConnector.Session)

	if a.sqsClient == nil {
		a.logger.Error(fmt.Sprintf("SQS client not created"))
		return
	}
}

func (a *AwsPingSender) SendPing(pingMessage string) error {
	if a.awsConnector == nil {
		return errors.New("AWS connector not active")
	}

	if a.awsConnector.Session == nil {
		return errors.New("AWS session not active")
	}

	if a.sqsClient == nil {
		return errors.New("AWS SQS client not active")
	}

	_, err := a.sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &a.queueUri,
		MessageBody: aws.String(pingMessage),
	})
	if err != nil {
		return err
	}

	return nil
}
