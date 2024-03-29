package resourceattrtocontextconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

// NewFactory creates a factory for the spanmetrics connector.
func NewFactory() connector.Factory {
	return connector.NewFactory(
		"resourceattr_to_context",
		createDefaultConfig,
		connector.WithTracesToTraces(createTracesToTraces, component.StabilityLevelDevelopment),
		connector.WithLogsToLogs(createLogsToLogs, component.StabilityLevelDevelopment),
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

func createLogsToLogs(
	_ context.Context,
	set connector.CreateSettings,
	cfg component.Config,
	logs consumer.Logs,
) (connector.Logs, error) {
	return newLogsConnector(set, cfg, logs)
}
