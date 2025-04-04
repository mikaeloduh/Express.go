// Package body_parser provides middleware functions for parsing request bodies in different formats.
package bodyparser

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

// JSONDecoder decodes JSON data from an io.Reader into the provided value.
// It serves as a decoder function that can be registered with the request object
// to automatically parse JSON request bodies.
//
// Parameters:
//   - r: The io.Reader containing the JSON data to be decoded
//   - v: The target value where the decoded data will be stored
//
// Returns:
//   - error: Any error encountered during the decoding process
func JSONDecoder(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// XMLDecoder decodes XML data from an io.Reader into the provided value.
// It serves as a decoder function that can be registered with the request object
// to automatically parse XML request bodies.
//
// Parameters:
//   - r: The io.Reader containing the XML data to be decoded
//   - v: The target value where the decoded data will be stored
//
// Returns:
//   - error: Any error encountered during the decoding process
func XMLDecoder(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
