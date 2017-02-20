package main

import (
	"flag"
	"time"
	"log"
	grpcclient "github.com/ru-rocker/gokit-playground/lorem-grpc/client"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/ru-rocker/gokit-playground/lorem-grpc"
	"fmt"
	"strconv"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()
	loremService := grpcclient.New(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "lorem":
		var requestType, minStr, maxStr string

		requestType, args = pop(args)
		minStr, args = pop(args)
		maxStr, args = pop(args)

		min, _ := strconv.Atoi(minStr)
		max, _ := strconv.Atoi(maxStr)
		lorem(ctx, loremService, requestType, min, max)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

func lorem(ctx context.Context, service lorem_grpc.Service, requestType string, min int, max int) {
	mesg, err := service.Lorem(ctx, requestType, min, max)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(mesg)
}