// Package body_parser provides middleware functions for parsing request bodies in different formats.
// It contains middleware implementations for automatically detecting and setting appropriate
// decoders based on the Content-Type header of incoming requests.
package body_parser

import (
	"strings"

	"github.com/mikaeloduh/expressgo/pkg/expressgo"
)

// JSONBodyParser is a middleware that sets the BodyParser to JSONDecoder.
// It automatically detects JSON content based on the Content-Type header and configures
// the request to use the appropriate JSON decoder for parsing the request body.
//
// Parameters:
//   - w: The response writer for the HTTP request
//   - r: The HTTP request object containing headers and body
//   - next: The next middleware function in the chain
//
// Returns:
//   - error: Always returns nil as this middleware doesn't produce errors
func JSONBodyParser(w *expressgo.ResponseWriter, r *expressgo.Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		r.SetDecoder(JSONDecoder)
	}

	next()

	return nil
}

// XMLBodyParser is a middleware that sets the BodyParser to XMLDecoder.
// It automatically detects XML content based on the Content-Type header and configures
// the request to use the appropriate XML decoder for parsing the request body.
//
// Parameters:
//   - w: The response writer for the HTTP request
//   - r: The HTTP request object containing headers and body
//   - next: The next middleware function in the chain
//
// Returns:
//   - error: Always returns nil as this middleware doesn't produce errors
func XMLBodyParser(w *expressgo.ResponseWriter, r *expressgo.Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/xml") {
		r.SetDecoder(XMLDecoder)
	}

	next()

	return nil
}
