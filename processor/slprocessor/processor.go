package slprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type SlPocessor struct {
	logger      *zap.Logger
	cfg         *ProcessorConfig
	pingHandler *PingSender
}

func newSlPocessor() *SlPocessor {
	return &SlPocessor{}
}

func (sp *SlPocessor) init(cfg *ProcessorConfig, logger *zap.Logger) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	sp.cfg = cfg
	sp.logger = logger
	sp.pingHandler = NewPingSender()

	awsPingSender := NewAwsPingSender()
	if err := sp.pingHandler.Init(logger, cfg.CacheTimeExpirationMin, cfg.CacheCleanupIntervalMin, &awsPingSender); err != nil {
		return err
	}
	return nil
}

func (sp *SlPocessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if sp.shouldIgnoreSpans(ctx) {
		td.ResourceSpans().RemoveIf(func(td ptrace.ResourceSpans) bool {
			return true
		})

		agentInstanceId := sp.getAgentInstanceId(ctx)
		sp.logger.Info(fmt.Sprintf("Spans ignored one time for: %s", agentInstanceId))

		return td, nil
	}

	customerIdInterface := ctx.Value(sp.cfg.ConsumerIdFieldName)

	if customerIdInterface == nil {
		sp.logger.Error(fmt.Sprintf("CustomerId marker not found"))
		return td, nil
	}

	customerId, ok := customerIdInterface.(string)
	if !ok {
		sp.logger.Error(fmt.Sprintf("CustomerId marker is not a string"))
		return td, nil
	}

	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		resource := rs.Resource()
		resource.Attributes().PutStr("customerId", customerId)
	}

	sp.ReportAgentAlive(ctx, customerId)

	return td, nil
}

func (sp *SlPocessor) getAgentInstanceId(ctx context.Context) string {
	agentInstanceIdInterface := ctx.Value(sp.cfg.HttpHeaderNameAgentId)

	if agentInstanceIdInterface == nil {
		return ""
	}

	agentInstanceId := agentInstanceIdInterface.(string)

	return agentInstanceId
}

func (sp *SlPocessor) shouldIgnoreSpans(ctx context.Context) bool {
	ignoreTracesInterface := ctx.Value("IgnoreSpan")

	if ignoreTracesInterface == nil {
		return false
	}

	ignoreTraces := ignoreTracesInterface.(string)

	if ignoreTraces == "true" {
		return true
	}

	return false
}

func (sp *SlPocessor) ReportAgentAlive(ctx context.Context, customerId string) {
	agentInstanceId := sp.getAgentInstanceId(ctx)

	if agentInstanceId == "" {
		sp.logger.Warn(fmt.Sprintf("Agent instance id not present for CustomerId: %s", customerId))
		return
	}

	err := sp.pingHandler.ReportPing(agentInstanceId, customerId)
	if err != nil {
		sp.logger.Error(fmt.Sprintf("%v", err))
	}

	sp.logger.Info(fmt.Sprintf("Agent instance id: %s", agentInstanceId))
}
