package payload

import (
	"contestive/api/apierror"
	"contestive/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type field string

func validFieldRune(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_'
}

func (e *field) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, (*string)(e)); err != nil {
		return err
	}

	for _, c := range *e {
		if !validFieldRune(c) {
			return apierror.NewHttpError(fmt.Sprintf("Unsuported field name %s", *e), http.StatusBadRequest)
		}
	}

	*e = field(strings.ToLower(string(*e)))

	return nil
}

type sortParam struct {
	Field field
	Asc   bool
}

func (e *sortParam) UnmarshalJSON(buf []byte) error {
	var asc string
	tmp := []interface{}{&e.Field, &asc}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	e.Asc = asc != "DESC"
	return nil
}

type rangeParam struct {
	From int
	To   int
}

func (e *rangeParam) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&e.From, &e.To}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	return nil
}

type filterValues []interface{}

type filterParam struct {
	Field        field
	FilterValues filterValues
}

type filterParams []filterParam

func (e *filterParams) UnmarshalJSON(buf []byte) error {
	params := (*e)[:0]
	tmp := map[field]interface{}{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	fmt.Println(string(buf))

	for k, v := range tmp {
		switch val := v.(type) {
		case []interface{}:

			if len(val) > 0 {
				if val2, ok := val[0].([]interface{}); ok {
					params = append(params, filterParam{k, val2})
					break
				}
			}
			params = append(params, filterParam{k, val})
		default:
			params = append(params, filterParam{k, filterValues{val}})
		}
	}

	*e = params

	return nil
}

// func (e *FilterParams) SqlAllowedFieldsNoWhere(allowedFields map[string]string) string {
// 	if e == nil || len(*e) == 0 {
// 		return ""
// 	}

// 	builder := &strings.Builder{}
// 	for k, v := range *e {
// 		if internalName, ok := allowedFields[k]; ok && len(v) > 0 {
// 			if builder.Len() > 0 {
// 				builder.WriteRune(',')
// 			}
// 			fmt.Fprintf(builder, `"%s" in (`, internalName)
// 			for i, id := range v {
// 				if i > 0 {
// 					builder.WriteRune(',')
// 				}
// 				fmt.Fprintf(builder, "%d", id)
// 			}
// 			fmt.Fprintf(builder, ") ")
// 		}
// 	}

// 	return builder.String()
// }

// func (e *FilterParams) SqlAllowedFields(allowedFields map[string]string) string {
// 	res := e.SqlAllowedFieldsNoWhere(allowedFields)
// 	if len(res) > 0 {
// 		return "WHERE " + res
// 	}
// 	return ""
// }

// func (e *FilterParams) SqlAllowedFieldsWithAnd(allowedFields map[string]string) string {
// 	res := e.SqlAllowedFieldsNoWhere(allowedFields)
// 	if len(res) > 0 {
// 		return "AND " + res
// 	}
// 	return ""
// }

func (s jsonParseResponser) ParseQuery(r *http.Request) entity.ListOptions {
	q := entity.ListOptions{}
	values := r.URL.Query()
	sortValues := values["sort"]
	if len(sortValues) != 0 {
		sp := sortParam{}
		err := json.Unmarshal([]byte(sortValues[0]), &sp)
		if err == nil {
			q.Sort.Field = string(sp.Field)
			q.Sort.Asc = sp.Asc
		}
	}

	rangeValues := values["range"]
	if len(rangeValues) != 0 {
		rp := rangeParam{}
		err := json.Unmarshal([]byte(rangeValues[0]), &rp)
		if err == nil {
			q.Range.From = rp.From
			q.Range.To = rp.To
		}
	}

	filterValues := values["filter"]
	if len(filterValues) != 0 {
		fp := filterParams{}
		err := json.Unmarshal([]byte(filterValues[0]), &fp)
		if err == nil {
			q.Filter = make([]entity.FilterParam, len(fp))
			for i, v := range fp {
				q.Filter[i].Field = string(v.Field)
				q.Filter[i].FilterValues = v.FilterValues
			}
		}
	}

	return q
}

// type contextKey string

// var queryContextKey contextKey = "query"

// // QueryMiddleware middleware that parses common query parameters and puts them into context
// func QueryMiddleware(next http.Handler) http.Handler {
// 	hfn := func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, queryContextKey, ParseQuery(r))
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// 	return http.HandlerFunc(hfn)
// }

// func QueryFromContext(ctx context.Context) QueryParams {
// 	q := ctx.Value(queryContextKey)
// 	query, _ := q.(QueryParams)
// 	return query
// }
