package mongo

import (
	"fmt"

	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) Get(key string, q storage.ImageQueryOptions) (img *domain.Image, err error) {
	col := m.db.C(imagesC)
	query := bson.M{"key": key, "owner": q.Owner}
	proj := imgQueryOptsToProjection(q)

	err = col.Find(query).Select(proj).One(img)
	return
}

func (m *mongo) List(q storage.ImageQueryOptions, opts *storage.ListOptions) (imgs []domain.Image, err error) {
	if opts == nil {
		opts = storage.NewListOptions()
	}

	col := m.db.C(imagesC)
	query := bson.M{"owner": q.Owner}
	proj := imgQueryOptsToProjection(q)
	sort := imgQueryOptsToSort(q)
	limit := opts.PerPage
	skip := opts.Page * opts.PerPage

	err = col.Find(query).Select(proj).Sort(sort...).Limit(limit).Skip(skip).All(&imgs)
	return
}

func (m *mongo) Insert(img domain.Image) error {
	// Make sure empty array is written to document instead of null
	if len(img.Builds) == 0 {
		img.Builds = []domain.ImageBuild{}
	}

	return m.db.C(imagesC).Insert(img)
}

func (m *mongo) InsertBuild(img domain.Image, build domain.ImageBuild) error {
	col := m.db.C(imagesC)
	query := bson.M{"owner": img.Owner, "key": img.Key}

	var builds struct {
		Builds []domain.ImageBuild `bson:"builds"`
	}

	err := col.Find(query).Select(bson.M{"builds": 1}).One(&builds)
	if err != nil {
		return err
	}

	// FIXME: Max number of build history is hardcoded
	// If number of build records exceeds the maximum remove the first recorded
	// build record
	if len(builds.Builds) > 100 {
		err := col.Update(query, bson.M{"$pop": 1})
		if err != nil {
			return err
		}
	}

	return col.Update(query, bson.M{
		"$push": bson.M{
			"$each":     []domain.ImageBuild{build},
			"$position": 0,
		},
	})
}

func (m *mongo) Update(owner, key string, img domain.Image) error {
	col := m.db.C(imagesC)
	query := bson.M{"owner": owner, "key": key}
	update := bson.M{"$set": img}
	return col.Update(query, update)
}

func (m *mongo) UpdateBuild(owner, imgKey string, build domain.ImageBuild) error {
	col := m.db.C(imagesC)
	query := bson.M{"owner": owner, "key": imgKey}
	update := bson.M{
		"$set": bson.M{
			"builds.0.finished_at": build.FinishedAt,
			"builds.0.status":      build.Status,
		},
	}

	return col.Update(query, update)
}

func (m *mongo) Delete(owner string, key string) error {
	col := m.db.C(imagesC)
	query := bson.M{"owner": owner, "key": key}
	return col.Remove(query)
}

func imgQueryOptsToProjection(q storage.ImageQueryOptions) bson.M {
	fields := bson.M{}
	if q.WithStatus {
		fields["last_build_at"] = 1
		fields["builds"] = bson.M{"$slice": 1}
	}

	if q.WithHistory {
		fields["builds"] = 1
	}

	return fields
}

func imgQueryOptsToSort(q storage.ImageQueryOptions) []string {
	var sort []string
	for field, dir := range q.SortBy {
		switch dir {
		case storage.Asc:
			sort = append(sort, field)
		case storage.Desc:
			sort = append(sort, fmt.Sprintf("-%s", field))
		}
	}

	return sort
}
