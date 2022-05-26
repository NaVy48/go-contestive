package postgresql

import (
	"contestive/entity"
	"fmt"
	"strings"
)

type sqlListOptions struct {
	o             entity.ListOptions
	allowedFields []string
	sb            strings.Builder
	filterVars    []interface{}
}

func SqlListOptions(o entity.ListOptions, allowedFields []string) sqlListOptions {
	return sqlListOptions{o: o, allowedFields: allowedFields}
}
func (e *sqlListOptions) rangeSql() string {
	r := e.o.Range

	if !r.Valid() {
		r.From = 0
		r.To = 10
	}
	return fmt.Sprintf(" OFFSET %d LIMIT %d ", r.From, r.To-r.From+1)
}

func (e *sqlListOptions) orderSql() string {
	s := e.o.Sort

	if len(e.allowedFields) == 0 {
		return ""
	}

	field := e.AllowedOrDefaultField(s.Field)

	if field == "" {
		return ""
	}

	asc := "DESC"
	if s.Asc {
		asc = "ASC"
	}

	return fmt.Sprintf(` ORDER BY %s %s `, field, asc)
}

func (e *sqlListOptions) filterSql() string {
	f := e.o.Filter

	if len(e.allowedFields) == 0 {
		return " 1 = 1 "
	}

	e.sb.Reset()
	e.filterVars = e.filterVars[:0]

	for _, fp := range f {
		e.singleFilter(fp)
	}

	if e.sb.Len() == 0 {
		e.sb.WriteString(" 1 = 1 ")
	}

	return e.sb.String()
}

func (e *sqlListOptions) singleFilter(fp entity.FilterParam) {
	if !e.AllowedField(fp.Field) {
		return
	}

	if e.sb.Len() > 0 {
		e.sb.WriteString(" AND ")
	}
	e.sb.WriteString(fp.Field)
	e.sb.WriteString(" IN (")
	for i := range fp.FilterValues {
		if i == 0 {
			e.sb.WriteString("?")
		} else {
			e.sb.WriteString(", ?")
		}
	}
	e.sb.WriteString(") ")
	e.filterVars = append(e.filterVars, fp.FilterValues...)
}

func (e *sqlListOptions) AllowedField(name string) bool {
	for _, f := range e.allowedFields {
		if f == name {
			return true
		}
	}

	return false
}

func (e *sqlListOptions) AllowedOrDefaultField(name string) string {
	for _, f := range e.allowedFields {
		if f == name {
			return f
		}
	}

	return e.allowedFields[0]
}
