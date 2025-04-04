package expressgo

import (
	"fmt"
	"net/http"
)

type Request struct {
	*http.Request
	decoder Decoder
}

func NewRequest(r *http.Request) *Request {
	return &Request{r, nil}
}

func (r *Request) SetDecoder(dec Decoder) {
	r.decoder = dec
}

// ParseBodyInto decodes the request body into the provided object
func (r *Request) ParseBodyInto(obj any) error {
	if r.decoder == nil {
		return fmt.Errorf("body parser not set, content type: %s", r.Header.Get("Content-Type"))
	}

	return r.decoder(r.Body, obj)
}
