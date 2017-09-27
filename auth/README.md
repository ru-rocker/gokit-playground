# Auth
This is a sample for authentication using JWT

### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### transport.go
Handling about encode and decode json

### register.go
Register service to consul

### security.go
Handling authentication, creating JWT then stored into Consul KV

### logging.go
Logging function is under this file

### Running Consul

    docker run --rm -p 8400:8400 -p 8500:8500 -p 8600:53/udp -h node1 progrium/consul -server -bootstrap -ui-dir /ui

### execute

    cd $GOPATH/src/github.com/ru-rocker/gokit-playground
    go run auth/auth.d/main.go -consul.addr localhost -consul.port 8500 -advertise.addr 192.168.1.103 -advertise.port 7002
