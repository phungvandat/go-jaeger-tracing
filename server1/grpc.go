package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/phungvandat/jaegertracing/userproto"
	"google.golang.org/grpc"
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

	s := grpc.NewServer()
	userproto.RegisterUserSvcServer(s, &userSvc{})

	log.Printf("listening GRPC: localhost%s\n", addr)

	return s.Serve(lis)
}
