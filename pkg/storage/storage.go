package storage

import (
	"errors"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Storage represents different kind operations on the database
type Storage interface {
	// Close will shutdown the connection to storage server if any exists
	Close() error

	// FindImageByKey finds an image by its key
	FindImageByKey(key string) (domain.Image, error)

	// FindAllImages finds all images belongs to a user
	FindAllImages(username string) ([]domain.Image, error)

	// FindImageTags finds all the image tags belongs to the image
	FindImageTags(imageKey string) ([]domain.ImageTag, error)

	// FindUser finds a user by its username
	FindUser(username string) (domain.User, error)

	// CreateImage creates an image record
	CreateImage(image domain.Image) error

	// UpdateImage updates the matching image record by imageKey
	UpdateImage(imageKey string, image domain.Image) error

	// DeleteImage deletes the matching image record
	DeleteImage(imageKey string) error

	// Creates an image record tag belongs to given image
	CreateImageTag(imageKey string, tag domain.ImageTag) error

	// Updates a matching image tag record by imageName and tagName
	UpdateImageTag(imageKey string, tagName string, tag domain.ImageTag) error

	// DeleteImageTag deletes the matching image tag
	DeleteImageTag(imageKey string, tagName string) error

	// CreateUser creates a user record
	CreateUser(user domain.User) error

	// UpdateUser updates a user with matching username
	UpdateUser(username string, user domain.User) error
}

var (
	ErrNotFound = errors.New("not found")
)
