package expressgo

import (
	"net/http"
)

// Encoder is a function that encodes an object into a writer
type Encoder func(http.ResponseWriter, any) error

type EncoderDecorator func(Encoder) Encoder
