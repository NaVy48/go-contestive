package payload

import (
	"bytes"
	"contestive/api/apierror"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

type test_dto struct {
	Payload string `json:"payload"`
}

func newPayloadHandler() jsonParseResponser {
	return jsonParseResponser{log.New(io.Discard, "", 0)}
}

func TestParseJSON_NormalPayload(t *testing.T) {

	is := is.New(t)

	payload := test_dto{"Some payload"}
	reqBody, err := json.Marshal(payload)
	is.NoErr(err)

	r := httptest.NewRequest(http.MethodGet, "/api/something", bytes.NewReader(reqBody))

	request := test_dto{}
	err = newPayloadHandler().ParseJSON(r, &request)
	is.NoErr(err)
	is.Equal(payload, request)
}

func TestParseJSON_NonPointerValue(t *testing.T) {
	is := is.New(t)

	payload := test_dto{"Some payload"}
	reqBody, err := json.Marshal(payload)
	is.NoErr(err)

	r := httptest.NewRequest(http.MethodGet, "/api/something", bytes.NewReader(reqBody))

	request := test_dto{}
	err = newPayloadHandler().ParseJSON(r, request)
	httpError, success := err.(apierror.HttpError)
	is.True(success)
	is.Equal(httpError.HttpStatus(), http.StatusInternalServerError)
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	is := is.New(t)

	r := httptest.NewRequest(http.MethodGet, "/api/something", strings.NewReader(`{payload:"somepayload"}`))

	request := test_dto{}
	err := newPayloadHandler().ParseJSON(r, request)
	httpError, success := err.(apierror.HttpError)
	is.True(success)
	is.Equal(httpError.HttpStatus(), http.StatusBadRequest)
}
