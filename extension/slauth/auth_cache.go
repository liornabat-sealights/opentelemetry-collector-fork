package slauth

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type AuthCache struct {
	cacheClient *cache.Cache
}

func NewAuthCache() *AuthCache {
	return &AuthCache{}
}

func (ac *AuthCache) init(cfg *ExtensionConfig) error {
	ac.cacheClient = cache.New(cfg.TimeExpirationMin, cfg.CleanupIntervalMin)
	return nil
}

func (ac *AuthCache) add(cacheData *AuthCacheData, timeBeforeExpiration time.Duration) (*AuthCacheData, error) {
	if err := cacheData.Validate(); err != nil {
		return nil, err
	}

	ac.cacheClient.Set(cacheData.token, cacheData, timeBeforeExpiration)
	return cacheData, nil
}

func (ac *AuthCache) get(token string) (*AuthCacheData, bool) {
	if cacheData, found := ac.cacheClient.Get(token); found {
		return cacheData.(*AuthCacheData), found
	}
	return nil, false
}

func (ac *AuthCache) exists(token string) bool {
	_, found := ac.cacheClient.Get(token)
	return found
}
