package resourceattrtocontextconnector

import (
	"context"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
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
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		tracesForThisResource := ptrace.NewTraces()
		traces.ResourceSpans().At(i).CopyTo(tracesForThisResource.ResourceSpans().AppendEmpty())

		resource := traces.ResourceSpans().At(i).Resource()
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
		c.tracesConsumer.ConsumeTraces(newCtx, tracesForThisResource)
	}

	return nil
}

func (*tracesConnector) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
