package models

// PaginateDto is a struct that represents the data that is required to paginate a list of users
type PaginateDto struct {
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
}
