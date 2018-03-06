package storage

import (
	"github.com/mobingilabs/pullr/pkg/domain"
)

// ImageQueryOptions is used for configuring how the image data should be fetched
type ImageQueryOptions struct {
	WithStatus  bool
	WithHistory bool
	SortBy      map[string]Direction
	Owner       string
}

// ImageQuerySetter is a function which modifies the query
type ImageQuerySetter func(opts *ImageQueryOptions)

// ImageRepository is where the image data can be queried and updated
type ImageRepository interface {
	// Get an image record from the repository with matching image key
	Get(key string, q ImageQueryOptions) (*domain.Image, error)
	// List images from the repository
	List(q ImageQueryOptions, opts *ListOptions) ([]domain.Image, error)
	// Insert an image record to repository
	Insert(img domain.Image) error
	// InsertBuild adds a new build record to matching image
	InsertBuild(img domain.Image, build domain.ImageBuild) error
	// Update an image record
	Update(owner, key string, img domain.Image) error
	// UpdateBuild updates the latest build record for matching image
	UpdateBuild(owner, imgKey string, build domain.ImageBuild) error
	// Delete an image record by its key
	Delete(owner, key string) error
}

// NewImageQueryOpts creates an ImageQueryOptions object with given query set
func NewImageQueryOpts(owner string, options ...ImageQuerySetter) *ImageQueryOptions {
	opts := &ImageQueryOptions{
		Owner:       owner,
		WithStatus:  false,
		WithHistory: false,
	}

	for _, setter := range options {
		setter(opts)
	}

	return opts
}

// ImgWithStatus is a query setter to include last build status with fetched
// image records
func ImgWithStatus() ImageQuerySetter {
	return func(q *ImageQueryOptions) {
		q.WithStatus = true
	}
}

// ImgWithHistory is a query setter to include build history with fetched
// image records
func ImgWithHistory() ImageQuerySetter {
	return func(q *ImageQueryOptions) {
		q.WithHistory = true
	}
}

// ImgFilterByOwner is a query setter to filter out images owned by other users
func ImgFilterByOwner(owner string) ImageQuerySetter {
	return func(q *ImageQueryOptions) {
		q.Owner = owner
	}
}

// ImgSortByLastBuild is a query setter to sort fetched image data by their
// last build time
func ImgSortByLastBuild(dir Direction) ImageQuerySetter {
	return func(q *ImageQueryOptions) {
		q.SortBy["last_build_at"] = dir
	}
}

// ImgSortByUpdate is a query setter to sort fetched image data by their update
// time
func ImgSortByUpdate(dir Direction) ImageQuerySetter {
	return func(q *ImageQueryOptions) {
		q.SortBy["updated_at"] = dir
	}
}
