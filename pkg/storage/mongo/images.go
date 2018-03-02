package mongo

import (
	"math"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) FindImageByKey(key string) (domain.Image, error) {
	col := m.db.C(imagesC)

	var image domain.Image
	query := bson.M{"key": key}
	err := col.Find(query).Select(bson.M{"history": 0}).One(&image)
	return image, toStorageErr(err)
}

func (m *mongo) FindAllImages(username string, listOpts *storage.ListOptions) ([]domain.Image, storage.Pagination, error) {
	col := m.db.C(imagesC)

	var pagination storage.Pagination

	query := bson.M{"owner": username}
	count, err := col.Find(query).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, pagination, err
	}

	page := listOpts.GetPage()
	perPage := listOpts.GetPerPage()
	if count > perPage {
		pagination.Last = int(math.Max(math.Ceil(float64(count)/float64(perPage)), 1)) - 1
	} else {
		pagination.Last = 0
	}

	if page < pagination.Last {
		pagination.Next = page + 1
	} else {
		pagination.Next = page
	}

	pagination.PerPage = perPage
	pagination.Current = page
	pagination.Total = count

	var images []domain.Image
	err = col.Find(query).Sort(optsToMongoSort(listOpts)).Select(bson.M{"history": 0}).Limit(perPage).Skip(perPage * page).All(&images)

	return images, pagination, toStorageErr(err)
}

func (m *mongo) FindAllImagesSince(username string, since time.Time) ([]domain.Image, error) {
	col := m.db.C(imagesC)

	var images []domain.Image
	query := bson.M{
		"owner": username,
		"$or": []bson.M{
			{"updated_at": bson.M{"$gt": since}},
			{"created_at": bson.M{"$gt": since}},
		},
	}
	err := col.Find(query).Sort("name").Select(bson.M{"history": 0}).All(&images)

	return images, err
}

func (m *mongo) CreateImage(image domain.Image) (string, error) {
	image.Key = domain.ImageKey(image.Repository)
	err := m.db.C(imagesC).Insert(image)
	return image.Key, toStorageErr(err)
}

func (m *mongo) UpdateImage(oldKey string, image domain.Image) (string, error) {
	newKey := domain.ImageKey(image.Repository)
	if newKey == "" {
		newKey = oldKey
	}

	image.Key = newKey

	// Don't update the history
	updateFields := bson.M{"$set": image}
	delete(updateFields, "history")

	err := m.db.C(imagesC).Update(bson.M{"key": oldKey}, updateFields)
	return newKey, toStorageErr(err)
}

func (m *mongo) DeleteImage(imageKey string) error {
	err := m.db.C(imagesC).Remove(bson.M{"key": imageKey})
	return toStorageErr(err)
}
