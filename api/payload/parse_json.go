package payload

import (
	"contestive/api/apierror"
	"encoding/json"
	"log"
	"net/http"
)

var (
	ErrNotSuppoertedContentType = apierror.NewHttpError("request content type is not supported", http.StatusBadRequest)
)

// ParseJSON parses request body r to value v as JSON. v must be a pointer. return apierror.HttpError if failed
func (s jsonParseResponser) ParseJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		log.Println(err)
		switch err.(type) {
		case *json.InvalidUnmarshalError:
			return apierror.NewHttpErrorWrap(err, "internal server error", http.StatusInternalServerError)
		default:
			return apierror.NewHttpErrorWrap(err, "bad Request", http.StatusBadRequest)
		}
	}

	return nil
}
