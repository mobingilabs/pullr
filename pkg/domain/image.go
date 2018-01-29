package domain

import (
	"fmt"
)

// Image represents a docker image
type Image struct {
	// Key is a unique string with combination of "Repository.Provider:Repository.Owner:Repository.Name"
	Key            string     `json:"key" bson:"key"`
	Name           string     `json:"name" bson:"name"`
	Owner          string     `json:"owner" bson:"owner"`
	Repository     Repository `json:"repository" bson:"repository"`
	DockerfilePath string     `json:"dockerfile_path" bson:"dockerfile_path"`
}

// ImageTag represents docker tags for an Image
type ImageTag struct {
	ImageKey string `json:"image_key" bson:"image_key"`
	RefType  string `json:"ref_type" bson:"ref_type"`
	RefTest  string `json:"ref_test" bson:"ref_test"`
	Name     string `json:"name" bson:"name"`
}

// ImageKey generates a unique image key from the repository
func ImageKey(repo Repository) string {
	return fmt.Sprintf("%s:%s:%s", repo.Provider, repo.Owner, repo.Name)
}
