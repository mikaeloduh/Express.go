package expressgo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// JSONBodyEncoder is a middleware that configures the response writer to use JSON encoding.
// It sets the Content-Type header to application/json if the client accepts JSON format.
//
// Parameters:
//   - w: The ResponseWriter to configure
//   - r: The incoming Request containing headers
//   - next: The next middleware function in the chain
//
// Returns:
//   - error: Always returns nil as this middleware doesn't produce errors
func JSONBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoderDecorator(JSONEncoderDecorator)

	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" || strings.HasPrefix(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
	}

	next()

	return nil
}

// JSONEncoderDecorator creates a decorator for the encoder chain that handles JSON encoding.
// It checks if the Content-Type is set to application/json and uses the JSONEncoder if it is.
// Otherwise, it passes the encoding task to the next encoder in the chain.
//
// Parameters:
//   - next: The next encoder in the chain to call if this one doesn't handle the encoding
//
// Returns:
//   - Encoder: A new encoder function that includes JSON encoding capability
func JSONEncoderDecorator(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj any) error {
		if w.Header().Get("Content-Type") == "application/json" {
			return JSONEncoder(w, obj)
		}
		return next(w, obj)
	}
}

// JSONEncoder encodes the provided object as JSON and writes it to the response writer.
// It uses the standard library's json package for encoding.
//
// Parameters:
//   - w: The http.ResponseWriter to write the encoded JSON to
//   - obj: The object to encode as JSON
//
// Returns:
//   - error: Any error that occurs during JSON encoding
func JSONEncoder(w http.ResponseWriter, obj any) error {
	return json.NewEncoder(w).Encode(obj)
}

// XMLBodyEncoder is a middleware that configures the response writer to use XML encoding.
// It sets the Content-Type header to application/xml if the client accepts XML format.
//
// Parameters:
//   - w: The ResponseWriter to configure
//   - r: The incoming Request containing headers
//   - next: The next middleware function in the chain
//
// Returns:
//   - error: Always returns nil as this middleware doesn't produce errors
func XMLBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoderDecorator(XMLEncoderDecorator)

	if strings.HasPrefix(r.Header.Get("Accept"), "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
	}

	next()

	return nil
}

// XMLEncoderDecorator creates a decorator for the encoder chain that handles XML encoding.
// It checks if the Content-Type is set to application/xml and uses the XMLEncoder if it is.
// Otherwise, it passes the encoding task to the next encoder in the chain.
//
// Parameters:
//   - next: The next encoder in the chain to call if this one doesn't handle the encoding
//
// Returns:
//   - Encoder: A new encoder function that includes XML encoding capability
func XMLEncoderDecorator(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj any) error {
		if w.Header().Get("Content-Type") == "application/xml" {
			return XMLEncoder(w, obj)
		}
		return next(w, obj)
	}
}

// XMLEncoder encodes the provided object as XML and writes it to the response writer.
// It uses the standard library's xml package for encoding.
//
// Parameters:
//   - w: The http.ResponseWriter to write the encoded XML to
//   - obj: The object to encode as XML
//
// Returns:
//   - error: Any error that occurs during XML encoding
func XMLEncoder(w http.ResponseWriter, obj any) error {
	return xml.NewEncoder(w).Encode(obj)
}
