package entity

type SortParam struct {
	Field string
	Asc   bool
}

type RangeParam struct {
	From int
	To   int
}

func (r *RangeParam) LimitAndDefault(limit, defaul int) {
	if r.To-r.From >= limit {
		r.To = r.From + limit - 1
	}
	if r.To-r.From <= 0 {
		r.To = r.From + defaul - 1
	}
}

func (r *RangeParam) Valid() bool {
	return r.To-r.From > 0
}

type FilterParam struct {
	Field        string
	FilterValues []interface{}
}

type ListOptions struct {
	Sort   SortParam
	Range  RangeParam
	Filter []FilterParam
}
