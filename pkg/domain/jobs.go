package domain

// BaseJob defines necessary information for all the jobs
type BaseJob struct {
	Source string `json:"source" valid:"required"`
	Action string `json:"action" valid:"required"`
}

// BuildImageJob describes necessary information to build a docker image
type BuildImageJob struct {
	BaseJob
	ImageKey   string `json:"image_key" valid:"required"`
	DockerTag  string `json:"tag" valid:"required"`
	CommitRef  string `json:"ref" valid:"required"`
	CommitHash string `json:"hash" valid:"required"`
}

// NewBuildImageJob creates a new job for building an image
func NewBuildImageJob(source, imageKey, ref, hash, dockerTag string) BuildImageJob {
	return BuildImageJob{
		BaseJob: BaseJob{
			Source: source,
			Action: "build",
		},

		ImageKey:   imageKey,
		CommitRef:  ref,
		DockerTag:  dockerTag,
		CommitHash: hash,
	}
}

type UpdateStatusJob struct {
	BaseJob
	Status
}

func NewUpdateStatusJob(source string, status Status) UpdateStatusJob {
	return UpdateStatusJob{
		BaseJob{Source: source, Action: "updatestatus"},
		status,
	}
}
