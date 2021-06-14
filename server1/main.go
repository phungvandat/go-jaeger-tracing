package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
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

	exitChn := make(chan error)
	go func() {
		osSignalChn := make(chan os.Signal, 1)
		signal.Notify(osSignalChn, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("exit by sign: %v\n", <-osSignalChn)
		exitChn <- nil
	}()

	go func() {
		err := grpcServer()
		if err != nil {
			exitChn <- err
		}
	}()

	go func() {
		err := httpServer()
		if err != nil {
			exitChn <- err
		}
	}()

	if err := <-exitChn; err != nil {
		log.Printf("server error: %v\n", err)
	}
	log.Printf("server stoped\n")
}
