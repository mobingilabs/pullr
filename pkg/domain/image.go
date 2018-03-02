package domain

import (
	"fmt"
	"time"
)

// Image status changing events
const (
	CauseImgBuildStart   = "images:build:start"
	CauseImgBuildFail    = "images:build:fail"
	CauseImgBuildSuccess = "images:build:success"
	CauseImgDelete       = "images:delete"
)

// Image statuses
const (
	StatusReady       = "images:ready"
	StatusImgBuilding = "images:building"
)

// Image represents a docker image
type Image struct {
	// Key is a unique string with combination of "Repository.Provider:Repository.Owner:Repository.Name"
	Key            string     `json:"key" bson:"key,omitempty"`
	Name           string     `json:"name" bson:"name,omitempty"`
	Owner          string     `json:"owner" bson:"owner,omitempty"`
	Repository     Repository `json:"repository" bson:"repository,omitempty"`
	DockerfilePath string     `json:"dockerfile_path" bson:"dockerfile_path,omitempty"`
	Tags           []ImageTag `json:"tags" bson:"tags,omitempty"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at,omitempty"`
}

// ImageTag represents docker tags for an Image
type ImageTag struct {
	RefType string `json:"ref_type" bson:"ref_type,omitempty"`
	RefTest string `json:"ref_test" bson:"ref_test,omitempty"`
	Name    string `json:"name" bson:"name,omitempty"`
}

// ImageKey generates a unique image key from the repository
func ImageKey(repo Repository) string {
	if repo.Name == "" || repo.Provider == "" || repo.Owner == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", repo.Provider, repo.Owner, repo.Name)
}
