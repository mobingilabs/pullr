package storage

import (
	"errors"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Storage represents different kind operations on the database
type Storage interface {
	// Close will shutdown the connection to storage server if any exists
	Close() error

	// FindImageByKey finds an image by its key
	FindImageByKey(key string) (domain.Image, error)

	// FindAllImagesSince finds images created after the given time
	FindAllImagesSince(s string, i time.Time) ([]domain.Image, error)

	// FindAllImages finds all images belongs to a user
	FindAllImages(username string) ([]domain.Image, error)

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

var (
	ErrNotFound = errors.New("not found")
)
