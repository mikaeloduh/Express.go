package framework

import (
	"fmt"
	"net/http"
)

type Request struct {
	*http.Request
	BodyParser Decoder
}

// ParseBodyInto decodes the request body into the provided object
func (r *Request) ParseBodyInto(obj interface{}) error {
	if r.BodyParser == nil {
		return fmt.Errorf("body parser not set, content type: %s", r.Header.Get("Content-Type"))
	}

	return r.BodyParser(r.Body, obj)
}
