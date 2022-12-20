package slauth

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/extension"
)

const (
	// The value of extension "type" in configuration.
	typeStr = "slauth"
)

// NewFactory creates a factory for the OIDC Authenticator extension.
// NewFactory creates a factory for FluentBit extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(typeStr, createDefaultConfig, createExtension, component.StabilityLevelBeta)
}

func createDefaultConfig() component.Config {
	return &ExtensionConfig{
		ExtensionSettings: config.NewExtensionSettings(component.NewID(typeStr)),
	}
}

func createExtension(_ context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	extCfg, ok := cfg.(*ExtensionConfig)
	if !ok {
		return nil, fmt.Errorf("extension config has unexpected type %T", cfg)
	}
	nse := NewSlauthExtension()
	if err := nse.init(extCfg, set.Logger); err != nil {
		return nil, err
	}

	return auth.NewServer(auth.WithServerStart(nse.start), auth.WithServerAuthenticate(nse.authenticate)), nil
}
