# lorem-hystrix
This is simple service module. Only for showing the micro service with HTTP and return json.
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

In this part I will demonstrate the circuit breaker patter. I copied from `lorem-consul`.

I am fully using all three functions from the golorem library.

## Required libraries

    go get github.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux
    go get github.com/juju/ratelimit
    go get github.com/prometheus/client_golang/prometheus
    go get github.com/go-kit/kit/sd/consul

### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### transport.go
Handling about encode and decode json

### logging.go
Logging function is under this file

### instrument.go
Middleware function. 
For this sample, this function for rate limiting and metrics.

### discovery.go
Consul service registration utility

#### lorem-consul.d
Go main function to build service and register to consul 

#### discover.d
Go main function for discover service

### Running Consul

    docker run --rm -p 8400:8400 -p 8500:8500 -p 8600:53/udp -h node1 progrium/consul -server -bootstrap -ui-dir /ui

### execute

    cd $GOPATH/src/github.com/ru-rocker/gokit-playground
    go run lorem-hystrix/lorem-consul.d/main.go -consul.addr localhost -consul.port 8500 -advertise.addr 192.168.1.103 -advertise.port 7002
    go run lorem-hystrix/discover.d/main.go -consul.addr localhost -consul.port 8500

### Running Prometheus and Grafana
To execute type

    cd $GOPATH/src/github.com/ru-rocker/gokit-playground
    docker-compose -f docker/docker-compose-prometheus-grafana-consul.yml up -d
    