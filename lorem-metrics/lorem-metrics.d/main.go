package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/ru-rocker/gokit-playground/lorem-metrics"
	ratelimitkit "github.com/go-kit/kit/ratelimit"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
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

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ru_rocker",
		Subsystem: "lorem_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "ru_rocker",
		Subsystem: "lorem_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc lorem_metrics.Service
	svc = lorem_metrics.LoremService{}
	svc = lorem_metrics.LoggingMiddleware(logger)(svc)
	svc = lorem_metrics.Metrics(requestCount, requestLatency)(svc)

	rlbucket := ratelimit.NewBucket(1*time.Second, 5)
	e := lorem_metrics.MakeLoremLoggingEndpoint(svc)
	e = ratelimitkit.NewTokenBucketThrottler(rlbucket, time.Sleep)(e)
	endpoint := lorem_metrics.Endpoints{
		LoremEndpoint: e,
	}

	r := lorem_metrics.MakeHttpHandler(ctx, endpoint, logger)

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
