package model

import "time"

type Content struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Created time.Time `json:"created"`
	Creator Creator   `json:"creator"`
}

type GetAllResp[T any] struct {
	Content       []T `json:"content"`
	TotalElements int `json:"totalElements"`
}

func (r GetAllResp[T]) Format() GetAllResp[T] {
	if r.Content == nil {
		r.Content = make([]T, 0, 0)
	}
	return r
}
