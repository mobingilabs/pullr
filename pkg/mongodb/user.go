package mongodb

import (
	"github.com/mobingilabs/pullr/pkg/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserStorage stores and queries user data from mongodb
type UserStorage struct {
	d *Driver
}

func (s *UserStorage) col() *mgo.Collection {
	return s.d.db.C(usersC)
}

// Get, gets a user by matching username
func (s *UserStorage) Get(username string) (domain.User, error) {
	var usr domain.User
	err := s.col().Find(bson.M{"username": username}).One(&usr)
	return usr, toStorageErr(err)
}

// GetByEmail, gets a user by matching email
func (s *UserStorage) GetByEmail(email string) (domain.User, error) {
	var usr domain.User
	err := s.col().Find(bson.M{"email": email}).One(&usr)
	return usr, toStorageErr(err)
}

// Put, puts a user record to mongodb
func (s *UserStorage) Put(user domain.User) error {
	err := s.col().Insert(user)
	return toStorageErr(err)
}

// Delete, deletes a user record by matching username
func (s *UserStorage) Delete(username string) error {
	err := s.col().Remove(bson.M{"username": username})
	return toStorageErr(err)
}
