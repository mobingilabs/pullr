package mongo

import (
	"github.com/mobingilabs/pullr/pkg/domain"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) FindUser(username string) (domain.User, error) {
	col := m.db.C(usersC)

	var user domain.User
	err := col.Find(bson.M{"username": username}).One(&user)
	return user, toStorageErr(err)
}

func (m *mongo) PutUserToken(username, provider, token string) error {
	usr, err := m.FindUser(username)
	if err != nil {
		return err
	}

	if usr.Tokens == nil {
		usr.Tokens = make(map[string]domain.UserToken)
	}

	usr.PutToken(provider, username, token)
	return m.UpdateUser(username, usr)
}

func (m *mongo) CreateUser(user domain.User) error {
	err := m.db.C(usersC).Insert(user)
	return toStorageErr(err)
}

func (m *mongo) UpdateUser(username string, user domain.User) error {
	err := m.db.C(usersC).Update(bson.M{"username": username}, user)
	return toStorageErr(err)
}
