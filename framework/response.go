package framework

import (
	"fmt"
	"net/http"
	"webframework/errors"
)

type ResponseWriter struct {
	http.ResponseWriter
	encoderHandler []EncoderHandler
}

func (w *ResponseWriter) UseEncoder(enc EncoderHandler) {
	w.encoderHandler = append(w.encoderHandler, enc)
}

func (w *ResponseWriter) Encode(obj interface{}) error {
	encoder := func(rw http.ResponseWriter, obj interface{}) error {
		return errors.NewError(http.StatusInternalServerError, fmt.Errorf("unsupported Content-Type: %s", rw.Header().Get("Content-Type")))
	}

	for i := len(w.encoderHandler) - 1; i >= 0; i-- {
		encoder = w.encoderHandler[i](encoder)
	}
	return encoder(w.ResponseWriter, obj)
}
