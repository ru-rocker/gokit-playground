package auth

import (
	"github.com/go-kit/kit/endpoint"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/leonelquinteros/gorand"
	"context"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("ru-rocker")
	method = jwt.SigningMethodHS256
)

func JwtEndpoint(consulAddress string, consulPort string, log log.Logger) endpoint.Middleware {

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(AuthRequest)
			response, err = next(ctx, request)
			resp := response.(AuthResponse)
			if strings.EqualFold("login", req.Type) {
				err = loginHandler(consulAddress, consulPort,
					req.Username, &resp, log)
			} else if strings.EqualFold("logout", req.Type) {
				println("logout")
				err = logoutHandler(consulAddress, consulPort,
					req, resp, log)
			}

			if err != nil {
				return nil, err
			}


			return resp, err
		}
	}
}
//create jwt keyFunc to retrieve kid
func keyFunc(token *jwt.Token) (interface{}, error) {
	return key, nil
}
// handling login
func loginHandler(consulAddress string, consulPort string,
	username string, resp *AuthResponse, log log.Logger) error {

	var (
		kid string
		tokenString string
	)

	defer func(){
		log.Log(
			"username", username,
			"kid", kid,
			"token", tokenString,
		)

	}()

	uuid, err := gorand.UUID()
	if err != nil {
		panic(err.Error())
	}
	kid = uuid

	claims := kitjwt.Claims{
		"user": username,
		"roles": resp.Roles,
	}

	token := jwt.NewWithClaims(method, jwt.MapClaims(claims))
	token.Header["kid"] = uuid

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(key)
	if err != nil {
		return err
	}

	resp.TokenString = tokenString
	log.Log(
		"username", username,
		"kid", uuid,
		"token", tokenString,
	)

	//register UUID on Consul KV
	client := ConsulClient(consulAddress, consulPort, log)
	kv := client.KV()
	p := &api.KVPair{Key: uuid, Value: []byte("active")}
	_, e := kv.Put(p, nil)
	if e != nil {
		return e
	}
	return nil
}

// handling logout
func logoutHandler(consulAddress string, consulPort string,
	req AuthRequest, resp AuthResponse, log log.Logger) error {

	var (
		username string
		kid string
		tokenString string
	)

	defer func(){
		log.Log(
			"username", username,
			"kid", kid,
			"token", tokenString,
		)

	}()

	tokenString = req.TokenString
	username = req.Username

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method != method {
			return nil, kitjwt.ErrUnexpectedSigningMethod
		}

		return keyFunc(token)
	})
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok && e.Inner != nil {
			if e.Errors&jwt.ValidationErrorMalformed != 0 {
				// Token is malformed
				return kitjwt.ErrTokenMalformed
			} else if e.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return kitjwt.ErrTokenExpired
			} else if e.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// Token is not active yet
				return kitjwt.ErrTokenNotActive
			}

			return e.Inner
		}

		return err
	}

	if !token.Valid {
		return kitjwt.ErrTokenInvalid
	}

	kid = token.Header["kid"].(string)

	//remove UUID on Consul KV
	client := ConsulClient(consulAddress, consulPort, log)
	kv := client.KV()
	_, e := kv.Delete (kid, nil)
	if e != nil {
		return e
	}

	return nil
}