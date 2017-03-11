package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"fmt"
	"github.com/ru-rocker/gokit-playground/lorem-logging"
	"github.com/go-kit/kit/log"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	logfile, err := os.OpenFile("./log/lorem/golorem.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()

	// Logging domain.
	var logger log.Logger
	{
		w := log.NewSyncWriter(logfile)
		logger = log.NewLogfmtLogger(w)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc lorem_logging.Service
	svc = lorem_logging.LoremService{}
	svc = lorem_logging.LoggingMiddleware(logger)(svc)
	endpoint := lorem_logging.Endpoints{
		LoremEndpoint: lorem_logging.MakeLoremLoggingEndpoint(svc),
	}

	r := lorem_logging.MakeHttpHandler(ctx, endpoint, logger)

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
