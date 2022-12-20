package slprocessor

import (
	"go.uber.org/zap"
)

type PingSenderInterface interface {
	Init(logger *zap.Logger, queueUri string)
	SendPing(pingMessage string) error
}
