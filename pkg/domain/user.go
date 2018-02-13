package domain

import (
	"golang.org/x/crypto/bcrypt"
)

// UserToken represents acquired oauth tokens
type UserToken struct {
	ID       string `json:"-" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Token    string `json:"token" bson:"password"`
}

// User defines both user authentication and relation
type User struct {
	Username       string               `json:"username" bson:"username,omitempty"`
	HashedPassword []byte               `json:"password,omitempty" bson:"password,omitempty"`
	Tokens         map[string]UserToken `json:"tokens" bson:"tokens,omitempty"`
}

// ComparePassword is the secure way to check if the given password is a match
// with user's password.
func (u User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	return err == nil
}

// PutToken appends the given pair of username and oauth token to user's token
// collection.
func (u *User) PutToken(provider, username, token string) {
	if u.Tokens == nil {
		u.Tokens = make(map[string]UserToken)
	}

	u.Tokens[provider] = UserToken{Username: username, Token: token}
}

// Token reports an oauth token by the given provider from user's token
// collection
func (u User) Token(provider string) *UserToken {
	if u.Tokens == nil {
		return nil
	}

	token, ok := u.Tokens[provider]
	if !ok {
		return nil
	}

	return &token
}
