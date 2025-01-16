package test

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webframework/framework"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Id       uint64 `json:"id" xml:"id"`
}

func TestReadBodyAsObject_JSON(t *testing.T) {
	body := TestObject{Username: "John Doe", Email: "jd@example.com", Id: 1}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/register", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	var testObject TestObject
	r := framework.Request{Request: req}
	r.BodyParser = framework.JSONDecoder

	err := r.ParseBodyInto(&testObject)
	assert.NoError(t, err)

	assert.Equal(t, body.Username, testObject.Username)
	assert.Equal(t, body.Email, testObject.Email)
	assert.Equal(t, body.Id, testObject.Id)
}

func TestReadBodyAsObject_XML(t *testing.T) {
	body := TestObject{Username: "John Doe", Email: "jd@example.com", Id: 1}
	xmlBody, _ := xml.Marshal(body)

	req := httptest.NewRequest("POST", "/register", strings.NewReader(string(xmlBody)))
	req.Header.Set("Content-Type", "application/xml")

	var testObject TestObject
	r := framework.Request{Request: req}
	r.BodyParser = framework.XMLDecoder

	err := r.ParseBodyInto(&testObject)
	assert.NoError(t, err)

	assert.Equal(t, body.Username, testObject.Username)
	assert.Equal(t, body.Email, testObject.Email)
	assert.Equal(t, body.Id, testObject.Id)
}

func TestReadBodyAsObject_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", strings.NewReader(""))
	req.Header.Set("Content-Type", "text/plain")

	var testObject TestObject
	r := framework.Request{Request: req}
	r.BodyParser = framework.JSONDecoder

	err := r.ParseBodyInto(&testObject)
	assert.Error(t, err)

	assert.Empty(t, testObject)
}

func TestReadBodyAsObject_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid body"))
	req.Header.Set("Content-Type", "application/json")

	var testObject TestObject
	r := framework.Request{Request: req}
	r.BodyParser = framework.JSONDecoder

	err := r.ParseBodyInto(&testObject)
	assert.Error(t, err)

	assert.ErrorContains(t, err, "invalid character 'i' looking for beginning of value")
	assert.Empty(t, testObject)
}

func TestWriteObjectAsJSON(t *testing.T) {
	testObject := TestObject{Username: "John Doe", Email: "jd@example.com", Id: 1}
	expected, _ := json.Marshal(testObject)

	wr := httptest.NewRecorder()
	w := &framework.ResponseWriter{ResponseWriter: wr}
	w.UseEncoder(framework.JSONEncoderHandler)
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
	w := &framework.ResponseWriter{ResponseWriter: wr}
	w.UseEncoder(framework.XMLEncoderHandler)
	w.Header().Set("Content-Type", "application/xml")

	err := w.Encode(testObject)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, wr.Code)
	assert.Equal(t, "application/xml", wr.Header().Get("Content-Type"))
	assert.Equal(t, string(expected), wr.Body.String())
}
