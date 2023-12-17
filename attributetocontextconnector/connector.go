package attributetocontextconnector

import (
	"context"

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
	honeycombApiKey, _ := traces.ResourceSpans().At(0).Resource().Attributes().Get("app.honeycomb_api_key")

	newCtx := context.WithValue(ctx, "x-honeycomb-team", honeycombApiKey)
	c.tracesConsumer.ConsumeTraces(newCtx, traces)
	return nil
}

func (*tracesConnector) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
