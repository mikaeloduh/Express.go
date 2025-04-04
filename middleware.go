package expressgo

// Middleware is a function that is called before the handler
type Middleware func(req *Request, res *Response, next func()) error
