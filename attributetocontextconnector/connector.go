package attributetocontextconnector

import (
	"context"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type tracesConnector struct {
	component.StartFunc
	component.ShutdownFunc

	logger *zap.Logger
	config *Config

	tracesConsumer consumer.Traces
}

func newTracesConnector(
	set connector.CreateSettings,
	config component.Config,
	traces consumer.Traces,
) (*tracesConnector, error) {
	return &tracesConnector{
		logger:         set.Logger,
		config:         config.(*Config),
		tracesConsumer: traces,
	}, nil
}

func (c *tracesConnector) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	honeycombApiKey, _ := traces.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes().Get("app.honeycomb_api_key")

	metadata := client.NewMetadata(map[string][]string{
		"x-honeycomb-team": {honeycombApiKey.AsString()},
	})
	clientInfo := client.Info{
		Metadata: metadata,
	}
	newCtx := client.NewContext(ctx, clientInfo)
	c.tracesConsumer.ConsumeTraces(newCtx, traces)
	return nil
}

func (*tracesConnector) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
