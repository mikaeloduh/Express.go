package expressgo

import (
	"net/http"
)

// Encoder is a function type interface that encodes an object into an HTTP response writer.
// It takes an http.ResponseWriter to write the encoded data to and any value to encode.
// Returns an error if the encoding process fails.
type Encoder func(http.ResponseWriter, any) error

// EncoderDecorator is a higher-order function interface that transforms an Encoder.
// It takes an existing Encoder and returns a new Encoder with enhanced functionality.
type EncoderDecorator func(Encoder) Encoder
