package expressgo

import (
	"net/http"
	"strings"

	"github.com/mikaeloduh/expressgo/e"
)

// Handler is a function that implements the Handler interface
type Handler interface {
	ServeHTTP(*Request, *Response) error
}

// HandlerFunc is a function that implements the Handler interface
type HandlerFunc func(*Request, *Response) error

func (f HandlerFunc) ServeHTTP(req *Request, res *Response) error {
	return f(req, res)
}

// WrapHandler Convert the standard http.Handler to a Handler that returns an error
func WrapHandler(h http.Handler) Handler {
	return HandlerFunc(func(req *Request, res *Response) error {
		h.ServeHTTP(res.ResponseWriter, req.Request)
		return nil
	})
}

type Router struct {
	routes        map[string]map[string]Handler
	middlewares   []Middleware
	errorHandlers []ErrorHandlerFunc
}

func NewRouter() *Router {
	r := &Router{
		routes:        make(map[string]map[string]Handler),
		errorHandlers: []ErrorHandlerFunc{},
	}
	// register default error handlers
	r.RegisterErrorHandler(DefaultFallbackErrorHandler)
	r.RegisterErrorHandler(DefaultUnauthorizedErrorHandler)
	r.RegisterErrorHandler(DefaultNotFoundErrorHandler)
	r.RegisterErrorHandler(DefaultMethodNotAllowedErrorHandler)
	return r
}

// RegisterErrorHandler register an error handler
func (rt *Router) RegisterErrorHandler(handlerFunc ErrorHandlerFunc) {
	// add at the beginning of the handler chain
	rt.errorHandlers = append([]ErrorHandlerFunc{handlerFunc}, rt.errorHandlers...)
}

// HandleError handles errors
func (rt *Router) HandleError(err error, req *Request, res *Response) {
	if len(rt.errorHandlers) == 0 {
		// use default error handlers if no error handlers
		rt.errorHandlers = []ErrorHandlerFunc{DefaultNotFoundErrorHandler, DefaultMethodNotAllowedErrorHandler}
	}

	var currentHandlerIndex = 0
	var next func(error)
	next = func(err error) {
		if currentHandlerIndex >= len(rt.errorHandlers) {
			// use default error handler if no error handlers
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		handler := rt.errorHandlers[currentHandlerIndex]
		currentHandlerIndex++
		handler(err, req, res, next)
	}

	next(err)
}

// Use adds middleware to the router
func (rt *Router) Use(middleware ...Middleware) {
	rt.middlewares = append(rt.middlewares, middleware...)
}

// Handle registers a new route with a matcher for the URL path and method
func (rt *Router) Handle(path string, method string, handler Handler) {
	path = strings.Trim(path, "/")
	if path == "" {
		path = "/"
	}
	if _, ok := rt.routes[path]; !ok {
		rt.routes[path] = make(map[string]Handler)
	}

	rt.routes[path][method] = handler
}

// ServeHTTP handles incoming HTTP requests and dispatches them to the registered handlers.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := NewRequest(r)
	res := NewResponse(w)

	path := strings.Trim(req.URL.Path, "/")
	if path == "" {
		path = "/"
	}
	method := req.Method

	// check full path
	if methodHandlers, ok := rt.routes[path]; ok {
		if h, ok := methodHandlers[method]; ok {
			handler := rt.applyMiddleware(h)
			if err := handler.ServeHTTP(req, res); err != nil {
				rt.HandleError(err, req, res)
			}
			return
		}
		// 405
		rt.HandleError(e.ErrorTypeMethodNotAllowed, req, res)
		return
	}

	// 404
	rt.HandleError(e.ErrorTypeNotFound, req, res)
}

func (rt *Router) applyMiddleware(handler Handler) Handler {
	h := handler
	for i := len(rt.middlewares) - 1; i >= 0; i-- {
		mw := rt.middlewares[i]
		currentHandler := h
		h = HandlerFunc(func(r *Request, w *Response) error {
			var err error
			next := func() {
				err = currentHandler.ServeHTTP(r, w)
			}
			if err := mw(r, w, next); err != nil {
				return err
			}
			return err
		})
	}
	return h
}
