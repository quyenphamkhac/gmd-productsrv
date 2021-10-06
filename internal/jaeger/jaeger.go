package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/quyenphamkhac/gmd-productsrv/config"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/uber/jaeger-lib/metrics"
)

func InitJeagerTracing(cfg *config.Config) (opentracing.Tracer, io.Closer, error) {
	jaegerCfg := jaegercfg.Configuration{
		ServiceName: cfg.Jeager.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jeager.LogSpans,
			LocalAgentHostPort: cfg.Jeager.Host,
		},
	}

	return jaegerCfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
