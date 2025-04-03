package expressgo

import (
	"fmt"
	"net/http"

	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
)

type ResponseWriter struct {
	http.ResponseWriter
	encoder Encoder
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		encoder: func(rw http.ResponseWriter, obj any) error {
			// fallback encoder
			return e.NewError(http.StatusInternalServerError, fmt.Errorf("unsupported Content-Type: %s", rw.Header().Get("Content-Type")))
		},
	}
}

func (w *ResponseWriter) UseEncoderDecorator(enc EncoderDecorator) {
	w.encoder = enc(w.encoder)
}

func (w *ResponseWriter) Encode(obj any) error {
	return w.encoder(w.ResponseWriter, obj)
}
