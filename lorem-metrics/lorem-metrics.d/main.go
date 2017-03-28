package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/ru-rocker/gokit-playground/lorem-rate-limit"
	"time"
	"github.com/juju/ratelimit"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc lorem_rate_limit.Service
	svc = lorem_rate_limit.LoremService{}
	svc = lorem_rate_limit.LoggingMiddleware(logger)(svc)

	rlbucket := ratelimit.NewBucket(1*time.Second, 5)
	e := lorem_rate_limit.MakeLoremLoggingEndpoint(svc)
	e = lorem_rate_limit.NewTokenBucketLimiter(rlbucket)(e)
	endpoint := lorem_rate_limit.Endpoints{
		LoremEndpoint: e,
	}

	r := lorem_rate_limit.MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()


	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<- errChan)
}
