package expressgo

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
)

// ErrorHandlerFunc is an interface of error handler
type ErrorHandlerFunc func(err error, w *ResponseWriter, r *Request, next func(error))

// DefaultNotFoundErrorHandler return 404 page not found with detail message
func DefaultNotFoundErrorHandler(err error, w *ResponseWriter, r *Request, next func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		if errors.Is(er, e.ErrorTypeNotFound) {
			w.WriteHeader(er.Code)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte(fmt.Sprintf("Cannot find the path \"%v\"", r.URL.Path)))
			return
		}
	}

	next(err)
}

// DefaultMethodNotAllowedErrorHandler return 405 method not allowed with detail message
func DefaultMethodNotAllowedErrorHandler(err error, w *ResponseWriter, r *Request, next func(error)) {
	var er *e.Error
	if errors.As(err, &er) {
		if errors.Is(er, e.ErrorTypeMethodNotAllowed) {
			w.WriteHeader(er.Code)
			path := strings.Trim(r.URL.Path, "/")
			if path == "" {
				path = "/"
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte(fmt.Sprintf("Method \"%v\" is not allowed on path \"%v\"", r.Method, path)))
			return
		}
	}

	next(err)
}

func DefaultUnauthorizedErrorHandler(err error, w *ResponseWriter, r *Request, next func(error)) {
	if er, ok := err.(*e.Error); ok {
		if er == e.ErrorTypeUnauthorized {
			w.WriteHeader(er.Code)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("401 unauthorized"))
			return
		}
	}

	next(err)
}

// DefaultFallbackErrorHandler catch all remaining errors
func DefaultFallbackErrorHandler(err error, w *ResponseWriter, r *Request, next func(error)) {
	if er, ok := err.(*e.Error); ok {
		w.WriteHeader(er.Code)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(er.Error()))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("500 internal server error"))
}
