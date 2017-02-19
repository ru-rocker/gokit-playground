# lorem-grpc
This is simple service module. Only for showing the micro service with gRPC protocol
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

I am fully using all three functions from the golorem library.

## Required libraries

    go get github.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux

# pb
Protocol buffer module. The place to create proto files.
Download protoc from [here](https://github.com/google/protobuf/releases)
Then execute `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
*Note: don't forget to add GOBIN on your PATH*

### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### transport.go
Handling about encode and decode json

### execute
TODO