package domain

import "testing"

func TestImageKey(t *testing.T) {
	key := ImageKey(Image{
		Repository: SourceRepository{
			Name:     "pullr",
			Owner:    "mobingilabs",
			Provider: "github",
		},
	})

	if key != "github:mobingilabs:pullr" {
		t.Fail()
	}
}
