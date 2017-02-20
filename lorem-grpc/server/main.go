package main

import (
	"net"
	"flag"
	"github.com/ru-rocker/gokit-playground/lorem-grpc"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

func main() {

	var (
		gRPCAddr = flag.String("grpc", ":8081",
			"gRPC listen address")
	)
	flag.Parse()
	ctx := context.Background()

	// init lorem service
	var svc lorem_grpc.Service
	svc = lorem_grpc.LoremService{}
	errChan := make(chan error)

	// creating Endpoints struct
	endpoints := lorem_grpc.Endpoints{
		LoremEndpoint: lorem_grpc.MakeLoremEndpoint(svc),
	}

	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := lorem_grpc.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterLoremServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<- errChan)
}