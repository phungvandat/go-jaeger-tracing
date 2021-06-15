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
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("r", r)
			span.LogKV("panic", r)
			ext.Error.Set(span, true)
		}
	}()

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
	span, _ := opentracing.StartSpanFromContext(ctx, "sv2_phase1")
	defer span.Finish()

	time.Sleep(time.Duration(100+rand.Intn(500-100)) * time.Millisecond)
	var x *int
	y := *x + 1
	fmt.Println(y)
}

func httpServer() error {
	addr := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT_2"))
	log.Printf("listening HTTP: localhost%s\n", addr)
	return http.ListenAndServe(addr, &httpServe{})
}
