package handler

import (
	"contestive/entity"
	"net/http"
)

type ParseResponser interface {
	ParseQuery(r *http.Request) entity.ListOptions
	ParseJSON(r *http.Request, v interface{}) error
	ResponseJSON(w http.ResponseWriter, r *http.Request, v interface{})
	SetTotalCount(w http.ResponseWriter, r *http.Request, total int)
}
