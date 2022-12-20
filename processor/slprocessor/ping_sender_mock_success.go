package slprocessor

import (
	"errors"

	"go.uber.org/zap"
)

type PingSenderMockSuccess struct {
	logger         *zap.Logger
	queueUri       string
	wasMessageSent bool
}

func (a *PingSenderMockSuccess) Init(logger *zap.Logger, queueUri string) {
	a.logger = logger
	a.queueUri = queueUri
}

func (a *PingSenderMockSuccess) SendPing(pingMessage string) error {
	if a.wasMessageSent {
		return errors.New("The message was already sent")
	}

	a.wasMessageSent = true
	return nil
}
