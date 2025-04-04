package body_parser

import (
	"encoding/json"
	"encoding/xml"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mikaeloduh/expressgo"
)

type TestRequest struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestRequest_ReadBodyAsObject(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"field1":"value1","field2":123}`))
	req.Header.Set("Content-Type", "application/json")

	r := expressgo.Request{Request: req}
	r.SetDecoder(JSONDecoder)

	var testReq TestRequest
	err := r.ParseBodyInto(&testReq)
	assert.NoError(t, err)

	assert.Equal(t, "value1", testReq.Field1)
	assert.Equal(t, 123, testReq.Field2)
}

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
	r := expressgo.Request{Request: req}
	r.SetDecoder(JSONDecoder)

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
	r := expressgo.Request{Request: req}
	r.SetDecoder(XMLDecoder)

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
	r := expressgo.Request{Request: req}
	r.SetDecoder(JSONDecoder)

	err := r.ParseBodyInto(&testObject)
	assert.Error(t, err)

	assert.Empty(t, testObject)
}

func TestReadBodyAsObject_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid body"))
	req.Header.Set("Content-Type", "application/json")

	var testObject TestObject
	r := expressgo.Request{Request: req}
	r.SetDecoder(JSONDecoder)

	err := r.ParseBodyInto(&testObject)
	assert.Error(t, err)

	assert.ErrorContains(t, err, "invalid character 'i' looking for beginning of value")
	assert.Empty(t, testObject)
}
