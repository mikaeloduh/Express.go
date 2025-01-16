package framework

import "strings"

// JSONBodyParser is a middleware that sets the BodyParser to JSONDecoder
func JSONBodyParser(w *ResponseWriter, r *Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		r.BodyParser = JSONDecoder
	}

	next()

	return nil
}

func JSONBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoder(JSONEncoderHandler)

	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" || strings.HasPrefix(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
	}

	next()

	return nil
}

// XMLBodyParser is a middleware that sets the BodyParser to XMLDecoder
func XMLBodyParser(w *ResponseWriter, r *Request, next func()) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/xml") {
		r.BodyParser = XMLDecoder
	}

	next()

	return nil
}

func XMLBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoder(XMLEncoderHandler)

	if strings.HasPrefix(r.Header.Get("Accept"), "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
	}

	next()

	return nil
}
