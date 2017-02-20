package client

import (
	"github.com/ru-rocker/gokit-playground/lorem-grpc"
	"github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) lorem_grpc.Service {
	var loremEndpoint = grpctransport.NewClient(
		conn, "Lorem", "Lorem",
		lorem_grpc.EncodeGRPCLoremRequest,
		lorem_grpc.DecodeGRPCLoremResponse,
		pb.LoremResponse{},
	).Endpoint()

	return lorem_grpc.Endpoints{
		LoremEndpoint:     loremEndpoint,
	}
}