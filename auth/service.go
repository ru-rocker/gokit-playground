package auth

import (
	"strings"
	"errors"
)

type Service interface {

	Login(username string, password string) (mesg string, roles []string, err error)

	Logout() string

	// health check
	HealthCheck() bool

}

type AuthService struct {

}

// create type that return function.
// this will be needed in main.go
type ServiceMiddleware func(Service) Service

var InvalidLoginErr = errors.New("Username or Password does not equal")

func (AuthService) Login(username string, password string) (mesg string, roles []string, err error) {
	if strings.EqualFold("admin", username) &&
		strings.EqualFold("password", password) {
		mesg, roles, err = "Login succeed", []string{"Admin", "User"}, nil
	} else {
		mesg, roles, err = "", nil, InvalidLoginErr
	}
	return
}

func (AuthService) Logout() string {
	return "Logout Succeed."
}

func (AuthService) HealthCheck() bool {
	//to make the health check always return true is a bad idea
	//however, I did this on purpose to make the sample run easier.
	//Ideally, it should check things such as amount of free disk space or free memory
	return true
}