package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/phungvandat/jaegertracing/userproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed by load .env %v \n", err)
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	closeFunc := initJaeger()
	defer closeFunc()
	ctx := context.Background()
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT_1")),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed connect to svc a by error: %v", err)
	}
	defer conn.Close()
	client := userproto.NewUserSvcClient(conn)

	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(make(http.Header)),
	)
	span := tracer.StartSpan("call_grpc_server_1", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	headers := make(http.Header)
	err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
	if err != nil {
		log.Fatalf("failed get user by error: %v", err)
	}
	for k, v := range headers {
		if len(v) != 1 {
			continue
		}
		ctx = metadata.AppendToOutgoingContext(ctx, k, v[0])
	}
	ext.PeerService.Set(span, serviceName)
	ext.PeerAddress.Set(span, "/userproto.UserSvc/GetUser")

	res, err := client.GetUser(ctx, &userproto.GetUserReq{
		Id: 123,
	})
	if err != nil {
		log.Fatalf("failed get user by error: %v", err)
	}
	fmt.Println(res)
}
