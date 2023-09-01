package model

type Pager struct {
	// 页码
	Page int64 `json:"page,omitempty"`
	// 页数
	PageSize int64 `json:"page_size,omitempty"`
	// 总行数
	TotalRows int64 `json:"total_rows,omitempty"`
}
