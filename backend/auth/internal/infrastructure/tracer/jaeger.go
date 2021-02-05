package tracer

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jagercfg "github.com/uber/jaeger-client-go/config"
)

func InitGlobalTracer(name string, agentAddr string) (io.Closer, error) {
	jCfg := jagercfg.Configuration{
		Disabled:    false,
		ServiceName: name,
		Sampler: &jagercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jagercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddr,
		},
	}

	tracer, closer, err := jCfg.NewTracer(
		jagercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
