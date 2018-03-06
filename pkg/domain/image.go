package domain

import (
	"fmt"
	"time"
)

// ImageBuildStatus defines the status of a build
type ImageBuildStatus string

// Valid build statuses
const (
	ImgBuildInProgress = "in_progress"
	ImgBuildSucceed    = "succeed"
	ImgBuildFailed     = "failed"
)

// Image status changing events
const (
	CauseImgBuildStart   = "image:build:start"
	CauseImgBuildFail    = "image:build:fail"
	CauseImgBuildSuccess = "image:build:success"
	CauseImgDelete       = "image:delete"
)

// Image statuses
const (
	StatusReady       = "image:ready"
	StatusImgBuilding = "image:building"
)

// Image represents a docker image
type Image struct {
	// Key is a unique string with combination of "Repository.Provider:Repository.Owner:Repository.Name"
	Key            string       `json:"key" bson:"key,omitempty"`
	Name           string       `json:"name" bson:"name,omitempty"`
	Owner          string       `json:"owner" bson:"owner,omitempty"`
	Repository     Repository   `json:"repository" bson:"repository,omitempty"`
	DockerfilePath string       `json:"dockerfile_path" bson:"dockerfile_path,omitempty"`
	Tags           []ImageTag   `json:"tags" bson:"tags,omitempty"`
	CreatedAt      time.Time    `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time    `json:"updated_at" bson:"updated_at,omitempty"`
	LastBuildAt    *time.Time   `json:"last_build_at,omitempty" bson:"last_build_at,omitempty"`
	Builds         []ImageBuild `json:"builds" bson:"builds,omitempty"`
	Status         *Status      `json:"status" bson:"-"`
}

// ImageTag represents docker tags for an Image
type ImageTag struct {
	RefType string `json:"ref_type" bson:"ref_type,omitempty"`
	RefTest string `json:"ref_test" bson:"ref_test,omitempty"`
	Name    string `json:"name" bson:"name,omitempty"`
}

// ImageBuild represents a build process
type ImageBuild struct {
	StartedAt  time.Time `json:"started_at" bson:"started_at"`
	FinishedAt time.Time `json:"finished_at" bson:"finished_at"`
	Status     string
}

// NewImageBuild creates an ImageBuild instance for a new build process
func NewImageBuild() ImageBuild {
	return ImageBuild{
		StartedAt:  time.Now(),
		FinishedAt: time.Time{},
		Status:     ImgBuildInProgress,
	}
}

// ImageKey generates a unique image key from the repository
func ImageKey(repo Repository) string {
	if repo.Name == "" || repo.Provider == "" || repo.Owner == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s:%s", repo.Provider, repo.Owner, repo.Name)
}
