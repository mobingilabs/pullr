package domain

// Repository has the information for source code repository.
type Repository struct {
	Provider string `json:"provider" bson:"provider"`
	Owner    string `json:"owner" bson:"owner"`
	Name     string `json:"name" bson:"name"`
}
