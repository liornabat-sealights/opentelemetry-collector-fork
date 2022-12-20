package slprocessor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

type PingSender struct {
	agentCache    *AgentCache
	logger        *zap.Logger
	awsPingSender PingSenderInterface
}

func NewPingSender() *PingSender {
	return &PingSender{}
}

func (ps *PingSender) Init(logger *zap.Logger, expirationTimeMin int, cleanupIntervalMin int, pingSender *PingSenderInterface) error {
	ps.logger = logger

	agentCache := NewAgentCache()

	ps.logger.Info(fmt.Sprintf("DEFAULT EXPIRATION: %d, DEFAULT CLEANUP=%d", expirationTimeMin, cleanupIntervalMin))

	if err := agentCache.init(expirationTimeMin, cleanupIntervalMin); err != nil {
		return err
	}

	var pingQueueUri string = os.Getenv("SL_PROCESSOR_COCKPIT_QUEUE_URI")

	ps.agentCache = agentCache
	ps.awsPingSender = *pingSender
	ps.awsPingSender.Init(ps.logger, pingQueueUri)

	return nil
}

func (ps *PingSender) ReportPing(agentInstanceId string, customerId string) error {
	_, wasFound := ps.agentCache.get(agentInstanceId)

	if !wasFound {
		err := ps.sendPing(agentInstanceId, customerId)
		if err != nil {
			return err
		}

		err = ps.addToCache(agentInstanceId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *PingSender) constructPingMessage(agentInstanceId string, customerId string) *PingMessageDTO {
	now := time.Now().Unix()

	evnt := []PingEvent{{
		Type:           1016,
		Data:           nil,
		Origin:         "server",
		UtcTimestampMs: now,
	}}

	pingMessage := &PingMessageDTO{
		Type:       "add-agent-event-from-backend",
		Version:    "1.0",
		CustomerId: customerId,
		AgentId:    agentInstanceId,
		AppName:    "",
		Created:    now,
		Events:     evnt,
	}

	return pingMessage
}

func (ps *PingSender) sendPing(agentInstanceId string, customerId string) error {
	ps.logger.Info(fmt.Sprintf("AWS SQS ping sending for: (%s, %s)", customerId, agentInstanceId))

	pingMessage := ps.constructPingMessage(agentInstanceId, customerId)

	pingMessageBytes, err := json.Marshal(pingMessage)
	if err != nil {
		return err
	}

	pingMessageStr := string(pingMessageBytes)

	err = ps.awsPingSender.SendPing(pingMessageStr)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PingSender) addToCache(agentInstanceId string) error {
	if _, err := ps.agentCache.add(NewAgentInstanceCacheData().
		SetAgentInstanceId(agentInstanceId)); err != nil {
		return err
	}

	return nil
}

func (ps *PingSender) PrintCache() {
	ps.logger.Info(fmt.Sprintf("Printing cache: "))

	for k := range ps.agentCache.GetAllItems() {
		v, t, _ := ps.agentCache.getWithExpiration(k)
		ps.logger.Info(fmt.Sprintf("Key=%s, value=%s, time=%s", k, v.agentInstanceId, t))
	}
}
