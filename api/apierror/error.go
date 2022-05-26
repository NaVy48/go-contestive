package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type HttpError struct {
	Message    string
	StatusCode int
	InnerError error
}

func (err HttpError) MarshalJSON() ([]byte, error) {
	dto := struct {
		Error      string `json:"error"`
		InnerError string `json:"innerError,omitempty"`
	}{Error: err.Message}

	if dev := os.Getenv("_DEVELOPMENT"); strings.ToLower(dev) == "true" && err.InnerError != nil {
		dto.InnerError = err.InnerError.Error()
	}

	return json.Marshal(dto)
}

func (err HttpError) HttpStatus() int {
	if err.StatusCode != 0 {
		return err.StatusCode
	}
	return http.StatusInternalServerError
}

func (err HttpError) Error() string {
	if err.InnerError != nil {
		return fmt.Sprintf("[HttpStatus: %d] %s: %s", err.StatusCode, err.Message, err.InnerError.Error())
	} else {
		return fmt.Sprintf("[HttpStatus: %d] %s", err.StatusCode, err.Message)
	}
}

func (e HttpError) Unwrap() error { return e.InnerError }

func NewHttpError(message string, statusCode int) HttpError {
	return HttpError{Message: message, StatusCode: statusCode}
}

func NewHttpErrorWrap(e error, message string, statusCode int) HttpError {
	return HttpError{Message: message, StatusCode: statusCode, InnerError: e}
}
