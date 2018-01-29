package domain

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User defines both user authentication and relation
type User struct {
	// Username is unique by user
	Username string `json:"username" bson:"username"`
	// Password is hash of the user's password
	Password []byte `json:"password,omitempty" bson:"password"`
}

func (u User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u User) ComparePassword(password string) bool {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Printf("%s == %s", u.Password, pass)
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	return err == nil
}
