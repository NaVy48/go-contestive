package payload

import (
	"contestive/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func compareFilter(got, want []entity.FilterParam, is *is.I) {
	for _, gotP := range got {
		found := false
		for _, wantP := range want {
			if gotP.Field == wantP.Field {
				found = true
				is.True(len(gotP.FilterValues) == len(wantP.FilterValues)) // Field has different length
				for i, v := range gotP.FilterValues {
					is.Equal(v, wantP.FilterValues[i])
				}
			}
		}
		is.True(found) // Filter field missing
	}
}

func Test_jsonParseResponser_ParseQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  entity.ListOptions
	}{
		// {"simple",
		// 	// filter={}&range=[0,9]&sort=["id","ASC"]
		// 	"filter=%7B%7D&range=%5B0%2C9%5D&sort=%5B%22id%22%2C%22ASC%22%5D",
		// 	entity.ListOptions{},
		// },
		{
			"empty query",
			"",
			entity.ListOptions{},
		},
		{
			"range query",
			// range=[0,9]
			"range=%5B0%2C9%5D",
			entity.ListOptions{Range: entity.RangeParam{From: 0, To: 9}},
		},
		{
			"sort query",
			// sort=["id","ASC"]
			"sort=%5B%22id%22%2C%22ASC%22%5D",
			entity.ListOptions{Sort: entity.SortParam{Field: "id", Asc: true}},
		},
		{
			"sort query missing asc",
			// sort=["id"]
			"sort=%5B%22id%22%5D",
			entity.ListOptions{Sort: entity.SortParam{Field: "id", Asc: true}},
		},
		{
			"sort query desc",
			// sort=["id","ASC"]
			"sort=%5B%22otherFieldName%22%2C%22DESC%22%5D",
			entity.ListOptions{Sort: entity.SortParam{Field: "otherFieldName", Asc: false}},
		},
		{
			"filter query",
			// filter={"id":[1,2,3],"name":"A"}
			"filter=%7B%22id%22:%5B1,2,3%5D,%22name%22:%22A%22%7D",
			entity.ListOptions{Filter: []entity.FilterParam{
				{Field: "id", FilterValues: []interface{}{1.0, 2.0, 3.0}},
				{Field: "name", FilterValues: []interface{}{"A"}},
			}},
		},
		{
			"filter query single id",
			// filter={"id":1}
			"filter=%7B%22id%22:1%7D",
			entity.ListOptions{Filter: []entity.FilterParam{
				{Field: "id", FilterValues: []interface{}{1.0}},
			}},
		},
		{
			"combined query",
			// filter={"id":1}&sort=["name","ASC"]&range=[10,20]
			"filter=%7B%22id%22:1%7D&sort=%5B%22name%22,%22ASC%22%5D&range=%5B10,20%5D",
			entity.ListOptions{
				Sort:  entity.SortParam{Field: "name", Asc: true},
				Range: entity.RangeParam{From: 10, To: 20},
				Filter: []entity.FilterParam{
					{Field: "id", FilterValues: []interface{}{1.0}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			s := newPayloadHandler()
			r := httptest.NewRequest(http.MethodGet, "http://example.com/t?"+tt.query, nil)
			got := s.ParseQuery(r)
			is.Equal(got.Range, tt.want.Range)
			is.Equal(got.Sort, tt.want.Sort)
			compareFilter(got.Filter, tt.want.Filter, is)
		})
	}
}
