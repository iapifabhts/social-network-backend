package model

type AllResp[T any] struct {
	Content       []T `json:"content"`
	TotalElements int `json:"totalElements"`
}

func (r AllResp[T]) Format() AllResp[T] {
	if r.Content == nil {
		r.Content = make([]T, 0, 0)
	}
	return r
}
