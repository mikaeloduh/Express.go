package expressgo

import (
	"fmt"
	"net/http"

	"github.com/mikaeloduh/expressgo/e"
)

type Response struct {
	http.ResponseWriter
	encoder Encoder
}

// NewResponse creates a new Response
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
		encoder: func(rw http.ResponseWriter, obj any) error {
			// fallback encoder
			return e.NewError(http.StatusInternalServerError, fmt.Errorf("unsupported Content-Type: %s", rw.Header().Get("Content-Type")))
		},
	}
}

func (rs *Response) UseEncoderDecorator(enc EncoderDecorator) {
	rs.encoder = enc(rs.encoder)
}

func (rs *Response) Encode(obj any) error {
	return rs.encoder(rs.ResponseWriter, obj)
}
