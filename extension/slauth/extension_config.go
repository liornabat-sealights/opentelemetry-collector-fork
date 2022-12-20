package slauth

import (
	"fmt"
	"go.opentelemetry.io/collector/config"
	"time"
)

type ExtensionConfig struct {
	config.ExtensionSettings `mapstructure:",squash"`
	IssuerURL                string        `mapstructure:"issuer_url"`
	IssuerUrlVerb            string        `mapstructure:"issuer_url_verb"`
	LogFileEnabled           bool          `mapstructure:"log_file_enabled"`
	TimeExpirationMin        time.Duration `mapstructure:"time_expiration_min"`
	CleanupIntervalMin       time.Duration `mapstructure:"cleanup_interval_min"`
	ConsumerIdFieldName      string        `mapstructure:"consumer_id_field_name"`
	HttpHeaderNameAgentId    string        `mapstructure:"agent_instance_id_field_name"`
}

func (c *ExtensionConfig) Validate() error {
	if c.IssuerURL == "" {
		return fmt.Errorf("issuer_url is required")
	}
	if c.IssuerUrlVerb == "" {
		return fmt.Errorf("issuer_url_verb is required")
	}
	if c.TimeExpirationMin <= 0 {
		return fmt.Errorf("time_expiration_min must be greater than 0")
	}
	if c.CleanupIntervalMin <= 0 {
		return fmt.Errorf("cleanup_interval_min must be greater than 0")
	}
	if c.ConsumerIdFieldName == "" {
		return fmt.Errorf("consumer_id_field_name is required")
	}

	return nil
}
