# lorem-logging
This is simple service module. Only for showing the micro service with HTTP and return json.
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

I am fully using all three functions from the golorem library.

## Required libraries

    go get github.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux


### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### transport.go
Handling about encode and decode json

### logging.go
Logging function is under this file

#### lorem-logging.d
Go main function will be located under this folder. The `dot d` means daemon.

### execute

    cd $GOPATH
    go run src/github.com/ru-rocker/gokit-playground/lorem/lorem.d/main.go

### Running Docker Command
docker run -v "/Users/ru-rocker/Documents/workspace-golang/src/github.com/ru-rocker/gokit-playground/lorem-logging/filebeat/filebeat.yml:/filebeat.yml" \
           -v "/Users/ru-rocker/golorem.log:/golorem.log" \
           -e "LOGSTASH_HOST=172.20.20.10" \
           -e "LOGSTASH_PORT=5044" \
           -e "INDEX=logstash" \
           fiunchinho/docker-filebeat
