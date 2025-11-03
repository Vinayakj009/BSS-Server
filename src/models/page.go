package models

type Page[V any] struct {
	TotalCount int64 `json:"total_count"`
	Items      []V   `json:"items"`
}

type PageableRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
