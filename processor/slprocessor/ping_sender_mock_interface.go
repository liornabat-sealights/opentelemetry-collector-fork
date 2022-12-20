package slprocessor

type PingSenderMockInterface interface {
	WasMessageSent() bool
	ClearSendStatus()
}
