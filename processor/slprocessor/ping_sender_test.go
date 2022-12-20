package slprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSenderCache_DoNotSendDuplicates(t *testing.T) {
	var awsSuccessfulMockSender PingSenderInterface = &PingSenderMockSuccess{}

	var ps *PingSender = &PingSender{}

	logger, _ := zap.NewProduction()

	ps.Init(logger, 1, 5, &awsSuccessfulMockSender)

	var wasSentTrial1 bool = false
	var wasSentTrial2 bool = false
	var wasSentTrial3 bool = false

	err := ps.ReportPing("123", "SeaLights")

	if err == nil {
		wasSentTrial1 = true
	}

	err = ps.ReportPing("123", "SeaLights")

	if err != nil {
		wasSentTrial2 = true
	}

	err = ps.ReportPing("123", "SeaLights")

	if err != nil {
		wasSentTrial3 = true
	}

	wasSentOnlyOnce := wasSentTrial1 && !wasSentTrial2 && !wasSentTrial3

	require.Equal(t, wasSentOnlyOnce, true)
}

func TestSenderCache_SendUniquePings(t *testing.T) {
	var awsSuccessfulMockSender PingSenderInterface = &PingSenderMockSuccess{}

	var ps *PingSender = &PingSender{}

	logger, _ := zap.NewProduction()

	ps.Init(logger, 1, 5, &awsSuccessfulMockSender)

	var wasSentTrial1 bool = false
	var wasSentTrial2 bool = false
	var wasSentTrial3 bool = false

	err := ps.ReportPing("123", "SeaLights")

	if err == nil {
		wasSentTrial1 = true
	}

	err = ps.ReportPing("456", "SeaLights")

	if err != nil {
		wasSentTrial2 = true
	}

	err = ps.ReportPing("789", "SeaLights")

	if err != nil {
		wasSentTrial3 = true
	}

	allMessagesWereSent := wasSentTrial1 && wasSentTrial2 && wasSentTrial3

	require.Equal(t, allMessagesWereSent, true)
}
