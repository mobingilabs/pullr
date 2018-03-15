package dummy

import (
	"sort"
	"strings"

	"github.com/mobingilabs/pullr/pkg/domain"
)

func sortUsers(users map[string]domain.User) []domain.User {
	sorted := make([]domain.User, 0, len(users))
	for name, usr := range users {
		index := sort.Search(len(sorted), func(i int) bool {
			return strings.Compare(sorted[i].Username, name) > 0
		})
		copy(sorted[index+1:], sorted[index:])
		sorted[index] = usr
	}

	return sorted
}

func sortImages(images map[string]domain.Image) []domain.Image {
	sorted := make([]domain.Image, 0, len(images))
	for name, img := range images {
		index := sort.Search(len(sorted), func(i int) bool {
			return strings.Compare(sorted[i].Name, name) >= 0
		})

		if index < len(sorted) {
			copy(sorted[index+1:], sorted[index:])
			sorted[index] = img
		} else {
			sorted = append(sorted, img)
		}
	}

	return sorted
}

func sortBuilds(builds []domain.Build) []domain.Build {
	sorted := make([]domain.Build, 0, len(builds))
	for _, build := range builds {
		index := sort.Search(len(sorted), func(i int) bool {
			return sorted[i].LastRecord.After(build.LastRecord)
		})
		copy(sorted[index+1:], sorted[index:])
		sorted[index] = build
	}

	return sorted
}

func sortImageBuilds(images map[string][]domain.Build) []domain.Build {
	sorted := make([]domain.Build, 0, len(images))
	for _, imgBuilds := range images {
		index := sort.Search(len(sorted), func(i int) bool {
			return sorted[i].LastRecord.After(imgBuilds[0].LastRecord)
		})
		copy(sorted[index+1:], sorted[index:])
		sorted[index] = imgBuilds[0]
	}

	return sorted
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
