package expressgo

// Middleware is a function that is called before the handler
type Middleware func(w *Response, r *Request, next func()) error
