package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	ilog "log"
	"flag"
	"context"
	"github.com/go-kit/kit/log"
	//kitjwt "github.com/go-kit/kit/auth/jwt"
	"net/http"
	"github.com/ru-rocker/gokit-playground/auth"
	//"github.com/dgrijalva/jwt-go"
)

func main() {
	// parse variable from input command
	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
		advertiseAddr = flag.String("advertise.addr", "", "advertise address")
		advertisePort = flag.String("advertise.port", "", "advertise port")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc auth.Service
	svc = auth.AuthService{}
	svc = auth.LoggingMiddleware(logger)(svc)

	e := auth.MakeAuthEndpoint(svc)
	e = auth.JwtEndpoint(*consulAddr, *consulPort, logger)(e)

	endpoint := auth.Endpoints{
		AuthEndpoint: e,
		HealthEndpoint: auth.MakeHealthEndpoint(svc),
	}

	r := auth.MakeHttpHandler(ctx, endpoint, logger)

	// Register Service to Consul
	registar := auth.Register(*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort)

	// HTTP transport
	go func() {
		ilog.Println("Starting server at port", *advertisePort)
		// register service
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe( ":" + *advertisePort, handler)
	}()


	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	error := <- errChan

	// deregister service
	registar.Deregister()
	ilog.Fatalln(error)
}