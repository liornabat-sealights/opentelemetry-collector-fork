package slprocessor

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type AwsConnector struct {
	logger  *zap.Logger
	Session *session.Session
}

func NewAwsConnector() *AwsConnector {
	return &AwsConnector{}
}

func (a *AwsConnector) awsEnvVariablesExists() bool {
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" &&
		os.Getenv("AWS_SECRET_ACCESS_KEY") != "" &&
		os.Getenv("AWS_REGION") != "" {
		return true
	}

	return false
}

func (a *AwsConnector) init(logger *zap.Logger) error {
	a.logger = logger

	if !a.awsEnvVariablesExists() {
		return errors.New("AWS environment variables not present")
	}

	var err error
	a.Session, err = session.NewSession()

	logger.Info(fmt.Sprintf("AWS session created"))

	if err != nil {
		return errors.New(fmt.Sprintf("AWS Session could not be made. Error: %v", err))
	}

	var r string = *(a.Session.Config.Region)
	logger.Info(fmt.Sprintf("AWS connection is operating in region: %s", r))

	return nil
}
