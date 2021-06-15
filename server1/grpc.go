package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/phungvandat/go-jaeger-tracing/userproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userSvc struct {
	userproto.UnimplementedUserSvcServer
}

func (s *userSvc) GetUser(ctx context.Context, req *userproto.GetUserReq) (*userproto.GetUserRes, error) {
	time.Sleep(200 * time.Millisecond)
	return &userproto.GetUserRes{
		User: &userproto.User{
			Id:       req.Id,
			Username: "example",
		},
	}, nil
}

func grpcServer() error {
	addr := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT_1"))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s\n", addr)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			mwUnary,
		),
	)
	userproto.RegisterUserSvcServer(s, &userSvc{})

	log.Printf("listening GRPC: localhost%s\n", addr)

	return s.Serve(lis)
}

func mwUnary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, span := mwTrace(ctx, info)
	defer span.Finish()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
			fmt.Println("r", r)
			span.LogKV("panic", r)
			ext.Error.Set(span, true)
		}
	}()

	return handler(ctx, req)
}

func mwTrace(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, opentracing.Span) {
	header, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		header = make(metadata.MD)
	}

	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header),
	)
	span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(spanCtx))
	ctx = opentracing.ContextWithSpan(ctx, span)

	return ctx, span
}
