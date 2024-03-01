package resourceattrtocontextconnector

import (
	"context"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type logsConnector struct {
	component.StartFunc
	component.ShutdownFunc

	logger *zap.Logger
	config *Config

	logsConsumer consumer.Logs
}

func newLogsConnector(
	set connector.CreateSettings,
	config component.Config,
	logs consumer.Logs,
) (*logsConnector, error) {
	return &logsConnector{
		logger:       set.Logger,
		config:       config.(*Config),
		logsConsumer: logs,
	}, nil
}

func (c *logsConnector) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		logsForThisResource := plog.NewLogs()
		logs.ResourceLogs().At(i).CopyTo(logsForThisResource.ResourceLogs().AppendEmpty())

		resource := logs.ResourceLogs().At(i).Resource()
		resourceAttributeMeta := make(map[string][]string)

		for j := 0; j < resource.Attributes().Len(); j++ {
			resource.Attributes().Range(func(key string, value pcommon.Value) bool {
				resourceAttributeMeta[key] = []string{value.AsString()}
				return true
			})
		}
		clientInfo := client.Info{
			Metadata: client.NewMetadata(resourceAttributeMeta),
		}
		newCtx := client.NewContext(ctx, clientInfo)
		c.logsConsumer.ConsumeLogs(newCtx, logsForThisResource)
	}

	return nil
}

func (*logsConnector) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
