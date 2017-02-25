package lorem_grpc

import (
	"github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
	"golang.org/x/net/context"
)

//Encode and Decode Lorem Request
func EncodeGRPCLoremRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(LoremRequest)
	return &pb.LoremRequest{
		RequestType: req.RequestType,
		Max: req.Max,
		Min: req.Min,
	} , nil
}

func DecodeGRPCLoremRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoremRequest)
	return LoremRequest{
		RequestType: req.RequestType,
		Max: req.Max,
		Min: req.Min,
	}, nil
}

// Encode and Decode Lorem Response
func EncodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(LoremResponse)
	return &pb.LoremResponse{
		Message: resp.Message,
		Err: resp.Err,
	}, nil
}

func DecodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.LoremResponse)
	return LoremResponse{
		Message: resp.Message,
		Err: resp.Err,
	}, nil
}