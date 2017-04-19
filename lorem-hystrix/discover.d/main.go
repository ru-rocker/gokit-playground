package main

import (
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/endpoint"
	"io"
	"strings"
	"net/url"
	"net/http"
	"context"
	ht "github.com/go-kit/kit/transport/http"
	consulsd "github.com/go-kit/kit/sd/consul"
	"os"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd/lb"
	"time"
	"github.com/gorilla/mux"
	"os/signal"
	"syscall"
	"fmt"
	"flag"
	"encoding/json"
	"strconv"
	"errors"
	"github.com/ru-rocker/gokit-playground/lorem-hystrix"
	"github.com/afex/hystrix-go/hystrix"
	"net"
)

//to execute: go run src/github.com/ru-rocker/gokit-playground/lorem-consul/discover.d/main.go -consul.addr localhost -consul.port 8500
// curl -XPOST -d'{"requestType":"word", "min":10, "max":10}' http://localhost:8080/sd-lorem
func main() {

	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger,"ts", log.DefaultTimestampUTC)
		logger = log.With(logger,"caller", log.DefaultCaller)
	}


	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()

		consulConfig.Address = "http://" + *consulAddr + ":" + *consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	tags := []string{"lorem", "ru-rocker"}
	passingOnly := true
	duration := 500 * time.Millisecond
	var loremEndpoint endpoint.Endpoint

	ctx := context.Background()
	r := mux.NewRouter()

	factory := loremFactory(ctx, "POST", "/lorem")
	serviceName := "lorem"
	subscriber := consulsd.NewSubscriber(client, factory, logger, serviceName, tags, passingOnly)
	balancer := lb.NewRoundRobin(subscriber)
	retry := lb.Retry(1, duration, balancer)
	loremEndpoint = retry

	// configure hystrix
	hystrix.ConfigureCommand("Lorem Request", hystrix.CommandConfig{Timeout: 1000})
	loremEndpoint = lorem_hystrix.Hystrix("Lorem Request", "Service currently unavailable", logger)(loremEndpoint)

	// POST /sd-lorem
	// Payload: {"requestType":"word", "min":10, "max":10}
	r.Methods("POST").Path("/sd-lorem").Handler(ht.NewServer(
		loremEndpoint,
		decodeConsulLoremRequest,
		lorem_hystrix.EncodeResponse, // use existing encode response since I did not change the logic on response
	))

	// Interrupt handler.
	errc := make(chan error)

	// configure the hystrix stream handler
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	}()
	
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", "8080")
		errc <- http.ListenAndServe(":8080", r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}

// factory function to parse URL from Consul to Endpoint
func loremFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc ht.EncodeRequestFunc
			dec ht.DecodeResponseFunc
		)
		enc, dec = encodeLoremRequest, decodeLoremResponse

		return ht.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}

// decode request from client (/sd-lorem)
// parsing JSON into LoremRequest
func decodeConsulLoremRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request lorem_hystrix.LoremRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// Encode request form LoremRequest into existing Lorem Service
// The encode will translate the LoremRequest into /lorem/{requestType}/{min}/{max}
func encodeLoremRequest(_ context.Context, req *http.Request, request interface{}) error {

	lr := request.(lorem_hystrix.LoremRequest)
	p := "/" + lr.RequestType + "/" + strconv.Itoa(lr.Min) + "/" + strconv.Itoa(lr.Max)
	req.URL.Path += p
	return nil
}

// decode response from Lorem Service
func decodeLoremResponse(_ context.Context, resp *http.Response) (interface{}, error) {

	var response lorem_hystrix.LoremResponse
	var s map[string]interface{}

	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil{
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}