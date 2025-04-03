package expressgo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Id       uint64 `json:"id" xml:"id"`
}

func TestWriteObjectAsJSON(t *testing.T) {
	testObject := TestObject{Username: "John Doe", Email: "jd@example.com", Id: 1}
	expected, _ := json.Marshal(testObject)

	wr := httptest.NewRecorder()
	w := &ResponseWriter{ResponseWriter: wr}
	w.UseEncoderDecorator(JSONEncoderDecorator)
	w.Header().Set("Content-Type", "application/json")

	err := w.Encode(testObject)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, wr.Code)
	assert.Equal(t, "application/json", wr.Header().Get("Content-Type"))
	assert.JSONEq(t, string(expected), wr.Body.String())
}

func TestWriteObjectAsXML(t *testing.T) {
	testObject := TestObject{Username: "John Doe", Email: "jd@example.com", Id: 1}
	expected, _ := xml.Marshal(testObject)

	wr := httptest.NewRecorder()
	w := &ResponseWriter{ResponseWriter: wr}
	w.UseEncoderDecorator(XMLEncoderDecorator)
	w.Header().Set("Content-Type", "application/xml")

	err := w.Encode(testObject)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, wr.Code)
	assert.Equal(t, "application/xml", wr.Header().Get("Content-Type"))
	assert.Equal(t, string(expected), wr.Body.String())
}
