package auth

import (
	"github.com/go-kit/kit/log"
	"time"
	"strings"
)

// implement function to return ServiceMiddleware
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

// Make a new type and wrap into Service interface
// Add logger property to this type
type loggingMiddleware struct {
	Service
	logger log.Logger
}

// Implement Service Interface for LoggingMiddleware
func (mw loggingMiddleware) Login(username string, password string) (mesg string, roles []string, err error)  {
	defer func(begin time.Time){
		mw.logger.Log(
			"function","Login",
			"mesg", mesg,
			"roles", strings.Join(roles, ","),
			"took", time.Since(begin),
		)
	}(time.Now())
	mesg, roles, err = mw.Service.Login(username, password)
	return
}

func (mw loggingMiddleware) Logout() (mesg string) {
	defer func(begin time.Time){
		mw.logger.Log(
			"function","Logout",
			"result", mesg,
			"took", time.Since(begin),
		)
	}(time.Now())
	mesg = mw.Service.Logout()
	return
}
