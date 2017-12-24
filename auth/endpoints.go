package auth

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"strings"
	"errors"
	"github.com/kr/pretty"
)

var ErrRequestTypeNotFound = errors.New("Request type only valid for login and logout")

type CommonReqResp struct{

	TokenString string `json:"-"`
}

//request
type AuthRequest struct {
	CommonReqResp
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"-"`
}

//response
type AuthResponse struct {
	CommonReqResp
	Roles []string `json:"roles,omitempty"`
	Mesg string `json:"mesg"`
	Err     error `json:"err,omitempty"`
}

//Health Request
type HealthRequest struct {

}

//Health Response
type HealthResponse struct {
	Status bool `json:"status"`
}

// endpoints wrapper
type Endpoints struct {
	AuthEndpoint endpoint.Endpoint
	HealthEndpoint endpoint.Endpoint
}

// creating Auth Endpoint
func MakeAuthEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			roles []string
			mesg string
			err error
		)

		req := request.(AuthRequest)
		pretty.Print("ctx")
		if strings.EqualFold(req.Type, "login") {
			mesg, roles, err = svc.Login(req.Username, req.Password)
		} else if strings.EqualFold(req.Type, "logout") {
			mesg = svc.Logout()
			err = nil
		} else {
			return nil, ErrRequestTypeNotFound
		}

		// check if err is not null
		if err != nil {
			return nil, err
		}
		return AuthResponse{Mesg:mesg, Roles: roles, Err: err}, nil
	}
}

// creating health endpoint
func MakeHealthEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status }, nil
	}
}
