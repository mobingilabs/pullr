package domain

const BuildQueue = "pullr-image-build"

type BaseImage struct {
	Source string `json:"source"`
	Action string `json:"action"`
}

type BuildImageJob struct {
	BaseImage
	ImageKey string `json:"image_key"`
}

func NewBuildImageJob(source, imageKey string) BuildImageJob {
	return BuildImageJob{
		BaseImage: BaseImage{
			Source: source,
			Action: "build",
		},

		ImageKey: imageKey,
	}
}
