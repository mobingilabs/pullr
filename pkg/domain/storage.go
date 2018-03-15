package domain

import "io"

// Directions
const (
	Asc  ListDir = "asc"
	Desc ListDir = "desc"
)

// StorageDriver is an interface wraps constructing various storage services
type StorageDriver interface {
	io.Closer
	AuthStorage() AuthStorage
	OAuthStorage() OAuthStorage
	UserStorage() UserStorage
	ImageStorage() ImageStorage
	BuildStorage() BuildStorage
}

// ListDir defines ordering/sorting direction
type ListDir string

// ListOptions is used for queries expected to report multiple records
type ListOptions struct {
	PerPage int `json:"per_page" query:"per_page"`
	Page    int `json:"page" query:"page"`
}

// Pagination contains pagination meta data about query
type Pagination struct {
	Last    int `json:"last"`
	Current int `json:"current"`
}

var DefaultListOptions = ListOptions{
	PerPage: 20,
	Page:    0,
}
