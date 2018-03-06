package domain

import (
	"time"
)

// Status represents a resource's status at a point in time. Metadata can be
// used for injecting additional information
type Status struct {
	Account  string    `json:"account" bson:"account"`
	Kind     string    `json:"kind" bson:"kind"`
	ID       string    `json:"id" bson:"id"`
	Time     time.Time `json:"time" bson:"time"`
	Name     string    `json:"name" bson:"name"`
	Metadata []byte    `json:"metadata" bson:"metadata"`
	Cause    string    `json:"cause" bson:"cause"`
}

const UnknownCause = "unknown"

func NewImageStatus(account string, imgKey string, name string, cause string, metadata []byte) Status {
	return Status{
		Kind:     "image",
		ID:       imgKey,
		Account:  account,
		Time:     time.Now(),
		Name:     name,
		Cause:    cause,
		Metadata: metadata,
	}
}
