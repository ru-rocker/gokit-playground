package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"fmt"
	"github.com/go-kit/kit/log"
	ratelimitkit "github.com/go-kit/kit/ratelimit"
	"github.com/ru-rocker/gokit-playground/lorem-rate-limit"
	"time"
	"golang.org/x/time/rate"
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

	limit := rate.NewLimiter(rate.Every(35*time.Millisecond), 100)
	e := lorem_rate_limit.MakeLoremLoggingEndpoint(svc)
	e = ratelimitkit.NewErroringLimiter(limit)(e)
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
