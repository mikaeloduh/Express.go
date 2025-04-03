package body_parser

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

func JSONDecoder(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

func XMLDecoder(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
