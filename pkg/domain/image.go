package domain

import (
	"fmt"
	"time"
)

// Image represents a docker image
type Image struct {
	// Key is a unique string with combination of "SourceRepository.Provider:SourceRepository.Owner:SourceRepository.Name"
	Key            string           `json:"key" bson:"key,omitempty"`
	Name           string           `json:"name" bson:"name,omitempty"`
	Owner          string           `json:"owner" bson:"owner,omitempty"`
	Repository     SourceRepository `json:"repository" bson:"repository,omitempty"`
	DockerfilePath string           `json:"dockerfile_path" bson:"dockerfile_path,omitempty"`
	Tags           []ImageTag       `json:"tags" bson:"tags,omitempty"`
	CreatedAt      time.Time        `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time        `json:"updated_at" bson:"updated_at,omitempty"`
}

// ImageTag represents docker tags for an Image
type ImageTag struct {
	RefType string `json:"ref_type" bson:"ref_type,omitempty"`
	RefTest string `json:"ref_test" bson:"ref_test,omitempty"`
	Name    string `json:"name" bson:"name,omitempty"`
}

// ImageStorage stores and queries image data
type ImageStorage interface {
	// Get retrieves a matching image record by username and key
	Get(username string, key string) (Image, error)

	// List retrieves a matching list of images
	List(username string, options ListOptions) ([]Image, Pagination, error)

	// Put inserts a new image record
	Put(username string, image Image) error

	// Finish updates a matching image record by username and key
	// with given image data
	Update(username string, key string, image Image) error

	// Delete deletes a matching image record
	Delete(username string, key string) error
}

// ImageKey generates a unique image key from the repository
func ImageKey(img Image) string {
	repo := img.Repository
	if repo.Name == "" || repo.Provider == "" || repo.Owner == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", repo.Provider, repo.Owner, repo.Name)
}
