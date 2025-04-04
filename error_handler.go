package expressgo

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mikaeloduh/expressgo/e"
)

// ErrorHandlerFunc is an interface of error handler
type ErrorHandlerFunc func(err error, req *Request, res *Response, next func(error))

// DefaultNotFoundErrorHandler return 404 page not found with detail message
func DefaultNotFoundErrorHandler(err error, req *Request, res *Response, next func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		if errors.Is(er, e.ErrorTypeNotFound) {
			res.WriteHeader(er.Code)
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = res.Write([]byte(fmt.Sprintf("Cannot find the path \"%v\"", req.URL.Path)))
			return
		}
	}

	next(err)
}

// DefaultMethodNotAllowedErrorHandler return 405 method not allowed with detail message
func DefaultMethodNotAllowedErrorHandler(err error, req *Request, res *Response, next func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		if errors.Is(er, e.ErrorTypeMethodNotAllowed) {
			res.WriteHeader(er.Code)
			path := strings.Trim(req.URL.Path, "/")
			if path == "" {
				path = "/"
			}
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = res.Write([]byte(fmt.Sprintf("Method \"%v\" is not allowed on path \"%v\"", req.Method, path)))
			return
		}
	}

	next(err)
}

func DefaultUnauthorizedErrorHandler(err error, _ *Request, res *Response, next func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		if errors.Is(er, e.ErrorTypeUnauthorized) {
			res.WriteHeader(er.Code)
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = res.Write([]byte("401 unauthorized"))
			return
		}
	}

	next(err)
}

// DefaultFallbackErrorHandler catch all remaining errors
func DefaultFallbackErrorHandler(err error, _ *Request, res *Response, _ func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		res.WriteHeader(er.Code)
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = res.Write([]byte(er.Error()))
		return
	}

	res.WriteHeader(http.StatusInternalServerError)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = res.Write([]byte("500 internal server error"))
}
