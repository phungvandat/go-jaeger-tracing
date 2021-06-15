package main

import (
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	serviceName = "server_1"
)

func initJaeger() func() {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // trace every call
		},
		Reporter: &config.ReporterConfig{
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
