package mongodb

import (
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersC   = "users"
	imagesC  = "images"
	historyC = "history"
)

// MongoDB is a `pullr.Storage` implementation for MongoDB
type MongoDB struct {
	Session *mgo.Session
	Db      *mgo.Database
}

// Dial connects to a MongoDB server and returns a pullr.Storage out of it
func Dial(servers string) (*MongoDB, error) {
	session, err := mgo.Dial(servers)
	if err != nil {
		return nil, err
	}

	mongodb := &MongoDB{
		Session: session,
		Db:      session.DB("pullr"),
	}

	return mongodb, nil
}

// Close closes the storage session, in this case connection to mongodb
func (s *MongoDB) Close() error {
	s.Session.Close()
	return nil
}

// FindImageByRepository implements Storage.FindImageByRepository
func (s *MongoDB) FindImageByKey(key string) (domain.Image, error) {
	col := s.Db.C(imagesC)

	var image domain.Image
	query := bson.M{"key": key}
	err := col.Find(query).One(&image)
	return image, toStorageErr(err)
}

// FindAllImages implements Storage.FindAllImages
func (s *MongoDB) FindAllImages(username string) ([]domain.Image, error) {
	col := s.Db.C(imagesC)

	var images []domain.Image
	query := bson.M{"owner": username}
	err := col.Find(query).Sort("-created_at").All(&images)

	return images, toStorageErr(err)
}

// FindAllImagesSince implements Storage.FindAllImagesSince
func (s *MongoDB) FindAllImagesSince(username string, since time.Time) ([]domain.Image, error) {
	col := s.Db.C(imagesC)

	var images []domain.Image
	query := bson.M{"owner": username, "created_at": bson.M{"$gt": since}}
	err := col.Find(query).Sort("-created_at").All(&images)

	return images, err
}

// FindUser finds a user records by uts username
func (s *MongoDB) FindUser(username string) (domain.User, error) {
	col := s.Db.C(usersC)

	var user domain.User
	err := col.Find(bson.M{"username": username}).One(&user)
	return user, toStorageErr(err)
}

// PutUserToken puts inserts given token to user's token list
func (s *MongoDB) PutUserToken(username, provider, token string) error {
	usr, err := s.FindUser(username)
	if err != nil {
		return err
	}

	if usr.Tokens == nil {
		usr.Tokens = make(map[string]string)
	}

	usr.Tokens[provider] = token
	return s.UpdateUser(username, usr)
}

// CreateImage implements Storage.CreateImage
func (s *MongoDB) CreateImage(image domain.Image) (string, error) {
	image.Key = domain.ImageKey(image.Repository)
	err := s.Db.C(imagesC).Insert(image)
	return image.Key, toStorageErr(err)
}

// UpdateImage implements Storage.UpdateImage
func (s *MongoDB) UpdateImage(oldKey string, image domain.Image) (string, error) {
	newKey := domain.ImageKey(image.Repository)
	if newKey == "" {
		newKey = oldKey
	}

	image.Key = newKey
	err := s.Db.C(imagesC).Update(bson.M{"key": oldKey}, bson.M{"$set": image})
	return newKey, toStorageErr(err)
}

// DeleteImage implements Storage.DeleteImage
func (s *MongoDB) DeleteImage(imageKey string) error {
	err := s.Db.C(imagesC).Remove(bson.M{"key": imageKey})
	return toStorageErr(err)
}

// CreateUser implements Storage.CreateUser
func (s *MongoDB) CreateUser(user domain.User) error {
	err := s.Db.C(usersC).Insert(user)
	return toStorageErr(err)
}

// UpdateUser implements Storage.UpdateUser
func (s *MongoDB) UpdateUser(username string, user domain.User) error {
	err := s.Db.C(usersC).Update(bson.M{"username": username}, user)
	return toStorageErr(err)
}

func toStorageErr(err error) error {
	switch err {
	case mgo.ErrNotFound:
		return storage.ErrNotFound
	default:
		return err
	}
}
