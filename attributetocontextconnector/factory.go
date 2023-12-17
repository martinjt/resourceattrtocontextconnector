package attributetocontextconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

// NewFactory creates a factory for the spanmetrics connector.
func NewFactory() connector.Factory {
	return connector.NewFactory(
		"honeycomb_attribute_to_context",
		createDefaultConfig,
		connector.WithTracesToTraces(createTracesToTraces, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesToTraces(
	_ context.Context,
	set connector.CreateSettings,
	cfg component.Config,
	traces consumer.Traces,
) (connector.Traces, error) {
	return newTracesConnector(set, cfg, traces)
}
