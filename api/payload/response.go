package payload

import (
	"bytes"
	"contestive/api/apierror"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var contentTypeKey = "Content-Type"
var contentLengthKey = "Content-Length"
var TotalCountKey = "X-Total-Count"
var jsonContetType = "application/json; charset=utf-8"

type HttpStatuser interface {
	HttpStatus() int
}

func (s jsonParseResponser) ResponseJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	responseJSON(w, r, s.log, v)
}

func responseJSON(w http.ResponseWriter, r *http.Request, log *log.Logger, v interface{}) {
	switch s := v.(type) {
	case apierror.HttpError:
		log.Println(s.Error())
		responseJSONWithStatus(w, r, log, v, s.HttpStatus())
	case error:
		responseJSON(w, r, log, apierror.NewHttpError(s.Error(), http.StatusInternalServerError))
	case HttpStatuser:
		responseJSONWithStatus(w, r, log, v, s.HttpStatus())
	default:
		responseJSONWithStatus(w, r, log, v, http.StatusOK)
	}
}

func responseJSONWithStatus(w http.ResponseWriter, r *http.Request, log *log.Logger, v interface{}, status int) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		responseJSON(w, r, log, apierror.NewHttpErrorWrap(err, "error while encoding to json", http.StatusInternalServerError))
		return
	}

	w.Header().Set(contentTypeKey, jsonContetType)
	w.Header().Set(contentLengthKey, strconv.Itoa(buf.Len()))
	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		// Should never occure
		log.Println("Error while writing the response", err)
	}
}

func (s jsonParseResponser) SetTotalCount(w http.ResponseWriter, r *http.Request, total int) {
	w.Header().Set(TotalCountKey, strconv.Itoa(total))
}
