package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type httpServe struct{}

func (s *httpServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	span := tracer.StartSpan(r.RequestURI, ext.RPCServerOption(spanCtx))
	defer span.Finish()

	span.Context().ForeachBaggageItem(func(k, v string) bool {
		fmt.Println(span, "baggage:", k, v)
		span.LogKV(k, v)
		return true
	})

	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	phase1(ctx)
	w.WriteHeader(http.StatusOK)
}

func phase1(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sv1_phase1")
	defer span.Finish()

	time.Sleep(time.Duration(100+rand.Intn(500-100)) * time.Millisecond)
	phase2(ctx)
}

func phase2(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sv1_phase2")
	defer span.Finish()

	time.Sleep(time.Duration(100+rand.Intn(500-100)) * time.Millisecond)
	phase3(ctx)
}

func phase3(ctx context.Context) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sv1_phase3")
	defer func() {
		defer span.Finish()
		if err == nil {
			return
		}
		span.LogKV("error", err)
		ext.Error.Set(span, true)
		log.Println("error: ", err)
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4321/call_svc_1", nil)
	if err != nil {
		return
	}

	span.LogKV("request", req.URL)
	ext.SpanKindRPCClient.Set(span)
	ext.PeerAddress.Set(span, req.RequestURI)
	ext.PeerService.Set(span, serviceName)
	ext.HTTPUrl.Set(span, req.RequestURI)
	ext.HTTPMethod.Set(span, http.MethodGet)

	tracer := opentracing.GlobalTracer()
	err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)) // inject Uber-Trace-Id to header
	if err != nil {
		return
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(req)
	if response != nil {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				log.Println("error close body", err)
			}
		}()
	}
	return nil
}

func httpServer() error {
	addr := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT_1"))
	log.Printf("listening HTTP: localhost%s\n", addr)
	return http.ListenAndServe(addr, &httpServe{})
}
