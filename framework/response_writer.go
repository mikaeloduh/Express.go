package framework

import "net/http"

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		encoderHandler: []EncoderHandler{},
	}
}
