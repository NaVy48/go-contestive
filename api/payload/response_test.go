package payload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestResponseJSON_NormalPayload(t *testing.T) {

	is := is.New(t)
	type dto struct {
		Payload string `json:"payload"`
	}

	payload := "Some payload"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/something", nil)
	newPayloadHandler().ResponseJSON(w, r, dto{payload})
	is.Equal(w.Code, http.StatusOK) // Status should be 200
	is.Equal(w.Result().Header.Get(contentTypeKey), jsonContetType)
	is.True(w.Result().ContentLength > 0)

	response := dto{}

	err := json.NewDecoder(w.Result().Body).Decode(&response)
	is.NoErr(err)
	is.Equal(payload, response.Payload)
}

func TestResponseJSON_CyclicPayload(t *testing.T) {

	is := is.New(t)
	type dto struct {
		Payload string `json:"payload"`
		Cycle   *dto   `json:"cycle"`
	}

	payload := dto{"Some payload", nil}
	payload.Cycle = &payload

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/something", nil)
	newPayloadHandler().ResponseJSON(w, r, payload)
	is.Equal(w.Code, http.StatusInternalServerError) // Status should be 500
}

type httpStatusPayload struct {
	Payload string `json:"payload"`
	status  int
}

func (p httpStatusPayload) HttpStatus() int {
	return p.status
}

func TestResponseJSON_HttpStatuser(t *testing.T) {
	is := is.New(t)

	payload := "Some payload"
	status := 555

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/something", nil)
	newPayloadHandler().ResponseJSON(w, r, httpStatusPayload{payload, status})
	is.Equal(w.Code, status) // Status should be 200
	is.Equal(w.Result().Header.Get(contentTypeKey), jsonContetType)
	is.True(w.Result().ContentLength > 0)

	response := httpStatusPayload{}

	err := json.NewDecoder(w.Result().Body).Decode(&response)
	is.NoErr(err)
	is.Equal(payload, response.Payload)
}

func TestResponseJSON_Error(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/something", nil)
	newPayloadHandler().ResponseJSON(w, r, fmt.Errorf("some error"))
	is.Equal(w.Code, http.StatusInternalServerError) // Status should be 500
}
