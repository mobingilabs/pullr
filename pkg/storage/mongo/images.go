package mongo

import (
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// List of common mongodb projections for images
var (
	imgsOnlyFirstBuild = bson.M{"builds": bson.M{"$slice": 1}}
	imgsNoHistory      = bson.M{"history": 0}
)

func (m *mongo) GetImages(keys []string) ([]domain.Image, error) {
	col := m.db.C(imagesC)

	var images []domain.Image
	query := bson.M{"key": bson.M{"$in": keys}}
	err := col.Find(query).All(&images)
	return images, err
}

func (m *mongo) FindImageByKey(key string) (domain.Image, error) {
	col := m.db.C(imagesC)

	var image domain.Image
	query := bson.M{"key": key}
	projection := mergeBson(imgsOnlyFirstBuild, imgsNoHistory)
	err := col.Find(query).Select(projection).One(&image)
	return image, toStorageErr(err)
}

func (m *mongo) FindImageByKeyWithBuilds(key string) (domain.Image, error) {
	col := m.db.C(imagesC)

	var image domain.Image
	query := bson.M{"key": key}
	projection := imgsNoHistory
	err := col.Find(query).Select(projection).One(&image)
	return image, toStorageErr(err)
}

func (m *mongo) FindAllImages(username string, listOpts *storage.ListOptions) ([]domain.Image, storage.Pagination, error) {
	col := m.db.C(imagesC)
	if listOpts == nil {
		listOpts = storage.NewListOptions()
	}

	query := bson.M{"owner": username}
	count, err := col.Find(query).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, storage.Pagination{}, err
	}

	pagination := storage.NewPagination(listOpts, count)
	limit := listOpts.PerPage
	skip := listOpts.PerPage * listOpts.Page
	projection := mergeBson(imgsOnlyFirstBuild, imgsNoHistory)

	var images []domain.Image
	err = col.Find(query).
		Select(projection).
		Sort(optsToMongoSort(listOpts)).
		Limit(limit).
		Skip(skip).
		All(&images)

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

	projection := mergeBson(imgsOnlyFirstBuild, imgsNoHistory)
	err := col.Find(query).Sort("name").Select(projection).All(&images)

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

func (m *mongo) StartImageBuild(username string, imgKey string, build domain.ImageBuild) error {
	var img domain.Image
	query := bson.M{"key": imgKey, "owner": username}
	err := m.db.C(imagesC).Find(query).One(&img)
	if err != nil {
		return err
	}

	// Truncate the array of builds if it exceeds maximum allowed number of builds
	// Since mongodb doesn't allow us to update the same field with two different
	// operators in the same time. We should pop the item first and then insert
	// the new value
	if len(img.Builds) >= 100 {
		update := bson.M{"$pop": bson.M{"builds": 1}}
		err := m.db.C(imagesC).Update(query, update)
		if err != nil {
			return err
		}
	}

	update := bson.M{
		"$set": bson.M{
			"last_build_at": time.Now(),
		},
		"$push": bson.M{
			"$each":     []domain.ImageBuild{build},
			"$position": 0,
		},
	}

	return m.db.C(imagesC).Update(query, update)
}

func (m *mongo) FinishImageBuild(username string, imgKey string, status domain.ImageBuildStatus) error {
	query := bson.M{"key": imgKey, "owner": username}
	update := bson.M{"builds.0.finished_at": time.Now(), "builds.0.status": status}
	return m.db.C(imagesC).Update(query, update)
}

func (m *mongo) DeleteImage(imageKey string) error {
	err := m.db.C(imagesC).Remove(bson.M{"key": imageKey})
	return toStorageErr(err)
}
