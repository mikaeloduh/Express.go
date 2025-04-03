package expressgo

import "io"

// Decoder is a function type that defines an interface for decoding data.
// It represents a function that takes an io.Reader and a destination object,
// then decodes the data from the reader into the provided object.
type Decoder func(io.Reader, any) error
