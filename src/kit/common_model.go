package kit

const (
	SortRequestMissingFirst = "first"
	SortRequestMissingLast  = "last"
)

type SortRequest struct {
	Field   string `json:"field"`
	Asc     bool   `json:"asc"`
	Missing string `json:"missing"`
}

type PagingRequest struct {
	Size   int            `json:"size"`
	Index  int            `json:"index"`
	SortBy []*SortRequest `json:"sortBy"`
}

type PagingResponse struct {
	Total int `json:"total"`
	Index int `json:"index"`
}
