package storage

import (
	"errors"
	"io"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Storage service errors
var (
	ErrNotFound = errors.New("not found")
)

// Service is responsible for persisting records
type Service interface {
	io.Closer
	// Users ===================================================================

	// FindUser finds a user by its username
	FindUser(username string) (domain.User, error)
	// PutUserToken add 3rd party service provider token to given user
	PutUserToken(username, provider, token string) error
	// CreateUser creates a user record
	CreateUser(user domain.User) error
	// UpdateUser updates a user with matching username
	UpdateUser(username string, user domain.User) error

	// Images ==================================================================

	// FindImageByKey finds an image by its key
	FindImageByKey(key string) (domain.Image, error)
	// FindAllImages finds all images belongs to a user, optional listOpts
	// parameter can be used to change ordering and limiting the data fetched
	// from the storage.
	FindAllImages(username string, listOpts *ListOptions) ([]domain.Image, Pagination, error)
	// FindAllImagesSince finds images created after the given time
	FindAllImagesSince(username string, since time.Time) ([]domain.Image, error)
	// CreateImage creates an image record and reports back its key
	CreateImage(image domain.Image) (string, error)
	// UpdateImage updates the matching image record by imageKey
	UpdateImage(imageKey string, image domain.Image) (string, error)
	// DeleteImage deletes the matching image record
	DeleteImage(imageKey string) error

	// History =================================================================

	// UpdateStatus updates the status of a given resource
	UpdateStatus(status domain.Status) error
	// Status reports back the last recorded status of a resource
	Status(username string, kind, id string) (*domain.Status, error)
	// Statuses reports back each resource's last status by the given resource
	// kind. Optional listOpts parameter can be used for pagination.
	// Sorting parameter of listOpts will be ignored and statuses will always
	// be sorted in chronological order.
	Statuses(username string, kind string, listOpts *ListOptions) ([]domain.Status, error)
	// StatusesByResources reports each resource's last status by the given
	// resource kind and matching id. Optional listOpts parameter can be used
	// for pagination. Sorting parameter will be ignored and statuses will
	// always be sorted in chronological order.
	StatusesByResources(username string, kind string, ids []string) ([]domain.Status, error)
	// StatusesByCause reports back each resource's last status for the given
	// resource kind. cause parameter is used for filtering the statuses by
	// their cause. Optional listOpts parameter can be used for pagination.
	// Sorting parameter of listOpts will be ignored and statuses will always
	// be sorted in chronological order.
	StatusesByCause(username, kind, cause string, listOpts *ListOptions) ([]domain.Status, error)
	// History reports all history for the given resource kind. Optional
	// listOpts parameter can be used for pagination. Sorting parameter of
	// listOpts will be ignored and statuses will always be sorted in
	// chronological order
	History(username, kind, id string, listOpts *ListOptions) ([]domain.Status, error)
}

// Pagination contains pagination meta data about query
type Pagination struct {
	Total   int `json:"total"`
	Next    int `json:"next_page"`
	Last    int `json:"last_page"`
	Current int `json:"current"`
	PerPage int `json:"per_page"`
}

// Direction defines ordering/sorting direction
type Direction string

// Directions
const (
	Asc  Direction = "asc"
	Desc           = "desc"
)

// ListOptions is used for queries expected to report multiple records
type ListOptions struct {
	PerPage       int               `query:"per_page"`
	Page          int               `query:"page"`
	SortBy        string            `query:"sort_by"`
	SortDirection Direction         `query:"sort_dir"`
	Filter        map[string]string `query:"filter"`
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
