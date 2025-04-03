package body_parser

import (
	"strings"

	"github.com/mikaeloduh/expressgo/pkg/expressgo"
)

// JSONBodyParser is a middleware that sets the BodyParser to JSONDecoder
func JSONBodyParser(w *expressgo.ResponseWriter, r *expressgo.Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		r.SetDecoder(JSONDecoder)
	}

	next()

	return nil
}

// XMLBodyParser is a middleware that sets the BodyParser to XMLDecoder
func XMLBodyParser(w *expressgo.ResponseWriter, r *expressgo.Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/xml") {
		r.SetDecoder(XMLDecoder)
	}

	next()

	return nil
}
