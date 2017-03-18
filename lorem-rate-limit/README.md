# lorem-rate-limit
This is simple service module. Only for showing the micro service with HTTP and return json.
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

In this part I will demonstrate how to limit a request based on Token Bucket Limiter algorithm.

I am fully using all three functions from the golorem library.

## Required libraries

    go get github.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux
    go get github.com/juju/ratelimit

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
For this sample, this function only for rate limiting only.

#### lorem-rate-limit.d
Go main function will be located under this folder. The `dot d` means daemon.

### execute

    cd $GOPATH/src/github.com/ru-rocker/gokit-playground
    go run lorem-logging/lorem-logging.d/main.go

### Running Filebeat
The filebeat is using docker-compose.
To execute type

    cd $GOPATH/src/github.com/ru-rocker/gokit-playground
    docker-compose -f docker/docker-compose-filebeat.yml up
    
*Notes: the log file is located under `$GOPATH/src/github.com/ru-rocker/gokit-playground/log/lorem/golorem.log`*