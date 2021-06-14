package main

import (
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func initJaeger() func() {
	cfg := config.Configuration{
		ServiceName: "service_a",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return func() {
		err := closer.Close()
		if err != nil {
			log.Printf("close tracer failed by error: %v", err)
		}
	}
}
