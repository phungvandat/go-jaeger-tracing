package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phungvandat/jaegertracing/userproto"
	"google.golang.org/grpc"
)

const (
	port = ":2345"
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

func main() {
	closeFunc := initJaeger()
	defer closeFunc()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %s\n", port)
	}

	s := grpc.NewServer()
	userproto.RegisterUserSvcServer(s, &userSvc{})
	log.Printf("listening localhost%s\n", port)

	exitChn := make(chan error)
	go func() {
		osSignalChn := make(chan os.Signal, 1)
		signal.Notify(osSignalChn, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("exit by sign: %v\n", <-osSignalChn)
		exitChn <- nil
	}()

	go func() {
		err := s.Serve(lis)
		if err != nil {
			exitChn <- err
		}
	}()

	if err := <-exitChn; err != nil {
		log.Printf("server error: %v\n", err)
	}
	log.Printf("server stoped\n")
}
