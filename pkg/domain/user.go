package domain

// User defines both user authentication and relation
type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

// UserToken represents acquired oauth tokens
type UserToken struct {
	ID       string `json:"-" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Token    string `json:"token" bson:"password"`
}

// UserStorage is the interface that handles persistence of authentication data
type UserStorage interface {
	Get(username string) (User, error)
	GetByEmail(username string) (User, error)
	Put(user User) error
	List(opts ListOptions) ([]User, Pagination, error)
	Delete(username string) error
}
