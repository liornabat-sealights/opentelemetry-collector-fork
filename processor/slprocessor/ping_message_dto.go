package slprocessor

type PingMessageDTO struct {
	Type       string      `json:"type"`
	Version    string      `json:"version"`
	CustomerId string      `json:"customerId"`
	AgentId    string      `json:"agentId"`
	AppName    string      `json:"appName"`
	Created    int64       `json:"created"`
	Events     []PingEvent `json:"events"`
}

type PingEvent struct {
	Type           int            `json:"type"`
	Data           *PingEventData `json:"data"`
	Origin         string         `json:"origin"`
	UtcTimestampMs int64          `json:"utcTimestamp_ms"`
}

type PingEventData struct{}
