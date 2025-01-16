package framework

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestRequest_ReadBodyAsObject(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", strings.NewReader(`{"field1":"value1","field2":123}`))
	req.Header.Set("Content-Type", "application/json")

	request := Request{Request: req, BodyParser: JSONDecoder}

	var testReq TestRequest
	err := request.ParseBodyInto(&testReq)
	assert.NoError(t, err)

	assert.Equal(t, "value1", testReq.Field1)
	assert.Equal(t, 123, testReq.Field2)
}
