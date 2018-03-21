package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mobingilabs/pullr/pkg/gova"
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

// Valid validates the image data
func (i Image) Valid() (bool, error) {
	validator := &gova.Validator{}
	validator.NotEmptyString("name", i.Name)
	validator.NotEmptyString("owner", i.Owner)
	validator.NotEmptyString("repository.provider", i.Repository.Provider)
	validator.NotEmptyString("repository.owner", i.Repository.Owner)
	validator.NotEmptyString("repository.name", i.Repository.Name)
	validator.NotEmptyString("dockerfile_path", i.DockerfilePath)
	validator.NotEmpty("tags", len(i.Tags))

	if len(i.Tags) > 0 {
		for index, tag := range i.Tags {
			_, errs := tag.Valid()
			validator.ExtendElt("tags", index, errs)
		}
	}

	return validator.Valid(), validator.Errors()
}

// MatchingTag reports back the matching build tag for given commit info
func (i Image) MatchingTag(commit *CommitInfo) (ImageTag, bool) {
	for _, tag := range i.Tags {
		if commit.RefType != tag.RefType {
			continue
		}

		if commit.RefType == SourceBranch {
			if commit.Ref == tag.RefTest {
				return tag, true
			}

			continue
		}

		test := tag.RefTest
		if test[0] == '/' && test[len(test)-1] == '/' {
			test = test[1 : len(test)-1]
		}

		match, err := regexp.MatchString(test, commit.Ref)
		if match && err == nil {
			return tag, true
		}
	}

	return ImageTag{}, false
}

// ImageTag represents docker tags for an Image
type ImageTag struct {
	RefType SourceRefType `json:"ref_type" bson:"ref_type,omitempty"`
	RefTest string        `json:"ref_test" bson:"ref_test,omitempty"`
	Name    string        `json:"name" bson:"name,omitempty"`
}

// Valid validates the image tag data
func (it ImageTag) Valid() (bool, gova.ValidationErrors) {
	validator := &gova.Validator{}
	validator.ShouldBeOneOf("ref_type", string(it.RefType), string(SourceTag), string(SourceBranch))
	validator.NotEmptyString("ref_test", it.RefTest)

	if it.RefType == SourceBranch {
		validator.NotEmptyString("name", it.Name)
	}

	return validator.Valid(), validator.Errors()
}

// Tag reports back container tag name
func (it ImageTag) Tag(commit *CommitInfo) string {
	if it.RefType == SourceBranch {
		return it.Name
	}

	if it.Name == "" {
		return commit.Ref
	}

	return it.Name
}

// ImageStorage stores and queries image data
type ImageStorage interface {
	// Get retrieves a matching image record belonging to user by its key
	Get(username string, key string) (Image, error)

	// GetMany retrieves matching image records belonging to user by their keys
	GetMany(username string, keys []string) (map[string]Image, error)

	// List retrieves a matching list of images
	List(username string, options ListOptions) ([]Image, Pagination, error)

	// Put inserts a new image record
	Put(username string, image Image) error

	// Update updates a matching image record by username and key with given image data
	Update(username string, key string, image Image) error

	// Delete deletes a matching image record
	Delete(username string, key string) error
}

// ImageKey generates a unique image key from the repository
func ImageKey(repo SourceRepository) string {
	if repo.Name == "" || repo.Provider == "" || repo.Owner == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", repo.Provider, repo.Owner, repo.Name)
}
