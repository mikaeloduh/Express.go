package expressgo

import "io"

// Decoder is a function interface that decodes a reader into an object
type Decoder func(io.Reader, any) error
