package middleware

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func JSONDecoder(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func XMLDecoder(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}
