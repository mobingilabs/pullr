package storage

import (
	"errors"
	"io"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Directions
const (
	Asc  Direction = "asc"
	Desc           = "desc"
)

// Storage errors
var (
	ErrNotFound = errors.New("not found")
)

// Storage represents different kind operations on the database
type Storage interface {
	io.Closer

	// FindImageByKey finds an image by its key
	FindImageByKey(key string) (domain.Image, error)

	// FindAllImagesSince finds images created after the given time
	FindAllImagesSince(s string, i time.Time) ([]domain.Image, error)

	// FindAllImages finds all images belongs to a user
	FindAllImages(username string, options *ListOptions) ([]domain.Image, Pagination, error)

	// FindUser finds a user by its username
	FindUser(username string) (domain.User, error)

	PutUserToken(username, provider, token string) error

	// CreateImage creates an image record and reports back its key
	CreateImage(image domain.Image) (string, error)

	// UpdateImage updates the matching image record by imageKey
	UpdateImage(imageKey string, image domain.Image) (string, error)

	// DeleteImage deletes the matching image record
	DeleteImage(imageKey string) error

	// CreateUser creates a user record
	CreateUser(user domain.User) error

	// UpdateUser updates a user with matching username
	UpdateUser(username string, user domain.User) error
}

// Direction defines ordering/sorting direction
type Direction string

// ListOptions is used for queries expected to report multiple records
type ListOptions struct {
	PerPage       int       `query:"per_page"`
	Page          int       `query:"page"`
	SortBy        string    `query:"sort_by"`
	SortDirection Direction `query:"sort_dir"`
}

// Pagination contains pagination meta data about query
type Pagination struct {
	Total   int `json:"total"`
	Next    int `json:"next_page"`
	Last    int `json:"last_page"`
	Current int `json:"current"`
	PerPage int `json:"per_page"`
}

// NewListOptions creates a new list options object
func NewListOptions(perPage, page int, sortBy string, sortDirection Direction) *ListOptions {
	return &ListOptions{perPage, page, sortBy, sortDirection}
}

// GetPerPage reports items per page for listing default value is 20
func (o *ListOptions) GetPerPage() int {
	if o == nil {
		return 20
	}

	return o.PerPage
}

// GetPage reports current page for listing default value is 0
func (o *ListOptions) GetPage() int {
	if o == nil {
		return 0
	}

	return o.Page
}

// GetSortBy reports sorting column for listing default is empty string
func (o *ListOptions) GetSortBy() string {
	if o == nil {
		return ""
	}

	return o.SortBy
}

// GetSortDirection reports sorting direction for listing default is ascending
func (o *ListOptions) GetSortDirection() Direction {
	if o == nil {
		return Asc
	}

	return o.SortDirection
}
