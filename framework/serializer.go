package framework

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

// Decoder is a function that decodes a reader into an object
type Decoder func(io.Reader, interface{}) error

func JSONDecoder(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func XMLDecoder(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}

// Encoder is a function that encodes an object into a writer
type Encoder func(http.ResponseWriter, interface{}) error

type EncoderHandler func(Encoder) Encoder

func JSONEncoderHandler(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj interface{}) error {
		if w.Header().Get("Content-Type") == "application/json" {
			return json.NewEncoder(w).Encode(obj)
		}
		return next(w, obj)
	}
}

func XMLEncoderHandler(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj interface{}) error {
		if w.Header().Get("Content-Type") == "application/xml" {
			return xml.NewEncoder(w).Encode(obj)
		}
		return next(w, obj)
	}
}
