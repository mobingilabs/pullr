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

func sortBuilds(records []domain.BuildRecord) []domain.BuildRecord {
	sorted := make([]domain.BuildRecord, 0, len(records))
	for _, record := range records {
		index := sort.Search(len(sorted), func(i int) bool {
			return sorted[i].StartedAt.After(record.StartedAt)
		})

		if index < len(sorted) {
			copy(sorted[index+1:], sorted[index:])
			sorted[index] = record
		} else {
			sorted = append(sorted, record)
		}
	}

	return sorted
}

func sortImageBuilds(images map[string]domain.Build) []domain.Build {
	sorted := make([]domain.Build, 0, len(images))
	for _, imgBuild := range images {
		index := sort.Search(len(sorted), func(i int) bool {
			return sorted[i].LastRecord.After(imgBuild.LastRecord)
		})

		if index < len(sorted) {
			copy(sorted[index+1:], sorted[index:])
			sorted[index] = imgBuild
		} else {
			sorted = append(sorted, imgBuild)
		}
	}

	return sorted
}
