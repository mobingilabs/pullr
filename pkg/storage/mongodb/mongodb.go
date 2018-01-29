package mongodb

import (
	"github.com/golang/glog"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersC     = "users"
	imagesC    = "images"
	imageTagsC = "imageTags"
	historyC   = "history"
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
	err := col.Find(query).All(&images)

	return images, toStorageErr(err)
}

// FindImageTags implements Storage.FindImageTags
func (s *MongoDB) FindImageTags(imageKey string) ([]domain.ImageTag, error) {
	col := s.Db.C(imageTagsC)

	var tags []domain.ImageTag
	err := col.Find(bson.M{"image_key": imageKey}).All(&tags)
	return tags, toStorageErr(err)
}

func (s *MongoDB) FindUser(username string) (domain.User, error) {
	col := s.Db.C(usersC)

	var user domain.User
	err := col.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		glog.Errorf("[ERROR] MongoDB.FindUser, %s", err)
	}
	return user, toStorageErr(err)
}

// CreateImage implements Storage.CreateImage
func (s *MongoDB) CreateImage(image domain.Image) error {
	image.Key = domain.ImageKey(image.Repository)
	err := s.Db.C(imagesC).Insert(image)
	return toStorageErr(err)
}

// UpdateImage implements Storage.UpdateImage
func (s *MongoDB) UpdateImage(oldKey string, image domain.Image) error {
	newKey := domain.ImageKey(image.Repository)
	image.Key = newKey
	if err := s.Db.C(imagesC).Update(bson.M{"key": oldKey}, image); err != nil {
		return toStorageErr(err)
	}

	if oldKey == newKey {
		return nil
	}

	// If image key is changed also update image_key fields of image tags
	tags, err := s.FindImageTags(oldKey)
	if err != nil {
		return toStorageErr(err)
	}

	for _, tag := range tags {
		tag.ImageKey = newKey
		if err := s.UpdateImageTag(oldKey, tag.Name, tag); err != nil {
			return toStorageErr(err)
		}
	}

	return toStorageErr(err)
}

// DeleteImage implements Storage.DeleteImage
func (s *MongoDB) DeleteImage(imageKey string) error {
	if err := s.Db.C(imagesC).Remove(bson.M{"key": imageKey}); err != nil {
		return err
	}

	err := s.Db.C(imageTagsC).Remove(bson.M{"image_key": imageKey})
	if err == mgo.ErrNotFound {
		return nil
	}

	return toStorageErr(err)
}

// CreateImageTag implements Storage.CreateImageTag
func (s *MongoDB) CreateImageTag(imageKey string, tag domain.ImageTag) error {
	tag.ImageKey = imageKey
	err := s.Db.C(imageTagsC).Insert(tag)
	return toStorageErr(err)
}

// UpdateImageTag implements Storage.UpdateImageTag
func (s *MongoDB) UpdateImageTag(imageKey string, tagName string, tag domain.ImageTag) error {
	err := s.Db.C(imageTagsC).Update(bson.M{"image_key": imageKey, "name": tagName}, tag)
	return toStorageErr(err)
}

// DeleteImageTag implements Storage.DeleteImageTag
func (s *MongoDB) DeleteImageTag(imageKey string, tagName string) error {
	err := s.Db.C(imageTagsC).Remove(bson.M{"image_key": imageKey, "name": tagName})
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
