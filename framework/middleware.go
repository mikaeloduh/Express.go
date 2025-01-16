package framework

// Middleware is a function that is called before the handler
type Middleware func(w *ResponseWriter, r *Request, next func()) error
