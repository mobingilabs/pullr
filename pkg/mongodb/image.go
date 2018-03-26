package mongodb

import (
	"github.com/mobingilabs/pullr/pkg/domain"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ImageStorage stores and queries image data from mongodb
type ImageStorage struct {
	d *Driver
}

func (s *ImageStorage) col() *mgo.Collection {
	return s.d.db.C(imagesC)
}

// Get finds an image by its owner and key
func (s *ImageStorage) Get(username string, key string) (domain.Image, error) {
	var image domain.Image
	query := bson.M{"key": key, "owner": username}
	err := s.col().Find(query).One(&image)
	return image, toStorageErr(err)
}

// GetMany find images by their keys which belongs to the given user
func (s *ImageStorage) GetMany(username string, keys []string) (map[string]domain.Image, error) {
	var images []domain.Image
	query := bson.M{"owner": username, "key": bson.M{"$in": keys}}
	err := s.col().Find(query).All(&images)
	if err != nil {
		return nil, toStorageErr(err)
	}

	imagesByKey := make(map[string]domain.Image, len(images))
	for _, img := range images {
		imagesByKey[img.Key] = img
	}

	return imagesByKey, toStorageErr(err)
}

// List reports back a list of images which belongs to a user.
func (s *ImageStorage) List(username string, opts domain.ListOptions) ([]domain.Image, domain.Pagination, error) {
	query := bson.M{"owner": username}
	count, err := s.col().Find(query).Count()
	if err != nil {
		return nil, domain.Pagination{}, toStorageErr(err)
	}

	skip, limit := opts.Cursor(count)
	pagination := opts.Paginate(count)

	var images []domain.Image
	err = s.col().Find(query).Skip(skip).Limit(limit).All(&images)
	return images, pagination, toStorageErr(err)
}

// Put puts an image record to mongodb database
func (s *ImageStorage) Put(image domain.Image) error {
	err := s.col().Insert(image)
	return toStorageErr(err)
}

// Update updates an image record in mongodb database
func (s *ImageStorage) Update(username string, key string, image domain.Image) error {
	err := s.col().Update(bson.M{"owner": username, "key": key}, image)
	return toStorageErr(err)
}

// Delete, deletes an image record by matching username and key
func (s *ImageStorage) Delete(username string, key string) error {
	err := s.col().Remove(bson.M{"owner": username, "key": key})
	return toStorageErr(err)
}
