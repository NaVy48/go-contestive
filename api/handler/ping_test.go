package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestPing(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	Ping().ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK) // Status should be 200
	reply := w.Body.String()
	is.Equal(reply, `pong`)
}
