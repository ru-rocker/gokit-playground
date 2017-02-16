# sentence
This is simple service module. Only for showing the micro service with HTTP and return json.
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

I am fully using all three functions from the golorem library.

## Required libraries

    go get githum.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux


### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### transport.go
Handling about encode and decode json

#### sentence.d
Go main function will be located under this folder. The `dot d` means daemon.

### execute

    cd $GOPATH
    go run src/github.com/ru-rocker/gokit-playground/player/player.d/main.go