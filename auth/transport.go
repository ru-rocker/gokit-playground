package auth

import (
	"github.com/gorilla/mux"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/kit/log"
	"context"
	"encoding/json"
	"net/http"
	"errors"
	"strings"
)

var ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")

// Make Http Handler
func MakeHttpHandler(_ context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	//POST /auth/{type}
	//type can be login or logout
	r.Methods("POST").Path("/auth/{type}").Handler(httptransport.NewServer(
		endpoint.AuthEndpoint,
		decodeAuthRequest,
		encodeResponse,
		options...,
	))

	//GET /health
	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoint.HealthEndpoint,
		decodeHealthRequest,
		encodeResponse,
		options...,
	))

	return r

}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err == InvalidLoginErr {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	println(err.Error())
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

// decode auth request
func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrBadRouting
	}

	var request AuthRequest
	if strings.EqualFold("login", requestType) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}
	}
	request.Type = requestType

	//get token from header
	val := r.Header.Get("Authorization")
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) == 2 && strings.ToLower(authHeaderParts[0]) == "bearer" {
		request.TokenString = authHeaderParts[1]
	}

	return request, nil
}

// encodeResponse is the common method to encode all response types to the
// client.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}

	if authResp, ok := response.(AuthResponse); ok {
		w.Header().Set("X-TOKEN-GEN", authResp.TokenString)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decode health check
func decodeHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}
