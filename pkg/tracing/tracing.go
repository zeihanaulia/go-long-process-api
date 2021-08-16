package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	samplerType  = "const"
	samplerParam = 1
)

func Init(serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("jaeger init error: %v", err)
	}

	cfg.ServiceName = serviceName
	cfg.Sampler.Type = samplerType
	cfg.Sampler.Param = samplerParam
	cfg.Reporter = &config.ReporterConfig{
		LogSpans: true,
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return tracer, closer, fmt.Errorf("jaeger init error: %v", err)
	}

	return tracer, closer, nil
}
