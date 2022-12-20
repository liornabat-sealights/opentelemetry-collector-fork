package slprocessor

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type AgentCache struct {
	cacheClient *cache.Cache
}

func NewAgentCache() *AgentCache {
	return &AgentCache{}
}

func (ac *AgentCache) init(defaultExpirationMin int, cleanupIntervalMin int) error {
	ac.cacheClient = cache.New(time.Duration(defaultExpirationMin)*time.Minute, time.Duration(cleanupIntervalMin)*time.Minute)
	return nil
}

func (ac *AgentCache) add(cacheData *AgentInstanceCacheData) (*AgentInstanceCacheData, error) {
	if err := cacheData.Validate(); err != nil {
		return nil, err
	}

	ac.cacheClient.Set(cacheData.agentInstanceId, cacheData, cache.DefaultExpiration)
	return cacheData, nil
}

func (ac *AgentCache) get(agentInstanceId string) (*AgentInstanceCacheData, bool) {
	if cacheData, found := ac.cacheClient.Get(agentInstanceId); found {
		return cacheData.(*AgentInstanceCacheData), found
	}

	return nil, false
}

func (ac *AgentCache) getWithExpiration(agentInstanceId string) (*AgentInstanceCacheData, time.Time, bool) {
	if cacheData, expiration, found := ac.cacheClient.GetWithExpiration(agentInstanceId); found {
		return cacheData.(*AgentInstanceCacheData), expiration, found
	}

	var expiration time.Time

	return nil, expiration, false
}

func (ac *AgentCache) exists(agentInstanceId string) bool {
	_, found := ac.cacheClient.Get(agentInstanceId)
	return found
}

func (ac *AgentCache) GetAllItems() map[string]cache.Item {
	return ac.cacheClient.Items()
}
