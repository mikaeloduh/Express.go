package expressgo

import (
	"fmt"
	"net/http"

	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
)

type ResponseWriter struct {
	http.ResponseWriter
	encoderHandler []EncoderHandler
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		encoderHandler: []EncoderHandler{},
	}
}

func (w *ResponseWriter) UseEncoder(enc EncoderHandler) {
	w.encoderHandler = append(w.encoderHandler, enc)
}

func (w *ResponseWriter) Encode(obj interface{}) error {
	encoder := func(rw http.ResponseWriter, obj interface{}) error {
		return e.NewError(http.StatusInternalServerError, fmt.Errorf("unsupported Content-Type: %s", rw.Header().Get("Content-Type")))
	}

	for i := len(w.encoderHandler) - 1; i >= 0; i-- {
		encoder = w.encoderHandler[i](encoder)
	}
	return encoder(w.ResponseWriter, obj)
}
