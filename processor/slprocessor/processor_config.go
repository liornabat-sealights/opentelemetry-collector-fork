package slprocessor

import (
	"fmt"

	"go.opentelemetry.io/collector/config"
)

type ProcessorConfig struct {
	config.ProcessorSettings `mapstructure:",squash"`
	ConsumerIdFieldName      string `mapstructure:"consumer_id_field_name"`
	HttpHeaderNameAgentId    string `mapstructure:"agent_instance_id_field_name"`
	CacheTimeExpirationMin   int    `mapstructure:"agent_cache_expiration_min"`
	CacheCleanupIntervalMin  int    `mapstructure:"agent_cache_cleanup_interval"`
}

func (c *ProcessorConfig) Validate() error {
	if c.ConsumerIdFieldName == "" {
		return fmt.Errorf("consumer_id_field_name is required")
	}

	return nil
}
