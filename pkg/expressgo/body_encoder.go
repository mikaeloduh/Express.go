package expressgo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

func JSONBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoderDecorator(JSONEncoderDecorator)

	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" || strings.HasPrefix(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
	}

	next()

	return nil
}

func JSONEncoderDecorator(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj any) error {
		if w.Header().Get("Content-Type") == "application/json" {
			return JSONEncoder(w, obj)
		}
		return next(w, obj)
	}
}

func JSONEncoder(w http.ResponseWriter, obj any) error {
	return json.NewEncoder(w).Encode(obj)
}

func XMLBodyEncoder(w *ResponseWriter, r *Request, next func()) error {
	w.UseEncoderDecorator(XMLEncoderDecorator)

	if strings.HasPrefix(r.Header.Get("Accept"), "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
	}

	next()

	return nil
}

func XMLEncoderDecorator(next Encoder) Encoder {
	return func(w http.ResponseWriter, obj any) error {
		if w.Header().Get("Content-Type") == "application/xml" {
			return XMLEncoder(w, obj)
		}
		return next(w, obj)
	}
}

func XMLEncoder(w http.ResponseWriter, obj any) error {
	return xml.NewEncoder(w).Encode(obj)
}
