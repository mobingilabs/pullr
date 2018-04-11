package domain

import (
	"io"
)

// Directions
const (
	Asc  ListDir = "asc"
	Desc ListDir = "desc"
)

// DefaultListOptions
var DefaultListOptions = ListOptions{
	PerPage: 20,
	Page:    0,
}

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

// Paginate creates a pagination info from list options
func (o ListOptions) Paginate(nitems int) Pagination {
	if o.PerPage == 0 {
		o.PerPage = 1
	}

	last := nitems / o.PerPage
	if last < 0 {
		last = 0
	}

	return Pagination{Last: last, Current: o.Page, Total: nitems}
}

// Cursor reports the cursor position for the current page and the number of
// items available in that page
func (o ListOptions) Cursor(nitems int) (int, int) {
	if nitems < o.PerPage {
		return 0, nitems
	}

	pos := o.PerPage * o.Page
	if pos > nitems {
		pos = nitems - o.PerPage
	}

	limit := o.PerPage
	if pos+o.PerPage > nitems {
		limit = nitems - pos
	}

	return pos, limit
}

// Pagination contains pagination meta data about query
type Pagination struct {
	Total   int `json:"total"`
	Last    int `json:"last"`
	Current int `json:"current"`
}
