package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phungvandat/jaegertracing/userproto"
	"google.golang.org/grpc"
)

func main() {
	closeFunc := initJaeger()
	defer closeFunc()

	conn, err := grpc.Dial("localhost:2345", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed connect to svc a by error: %v", err)
	}
	defer conn.Close()

	client := userproto.NewUserSvcClient(conn)
	ctx := context.Background()
	res, err := client.GetUser(ctx, &userproto.GetUserReq{
		Id: 123,
	})
	if err != nil {
		log.Fatalf("failed get user by error: %v", err)
	}
	fmt.Println(res)
}
