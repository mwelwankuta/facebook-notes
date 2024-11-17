package models

// GetPaginationFromQueryDto is a struct that represents the raw query parameters
type GetPaginationFromQueryDto struct {
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
}

// PaginateDto is a struct that represents the parsed pagination parameters
type PaginateDto struct {
	Offset int
	Limit  int
}
