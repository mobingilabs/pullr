package mongodb

import (
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// BuildStorage stores and queries build data from mongodb
type BuildStorage struct {
	d *Driver
}

func (s *BuildStorage) col() *mgo.Collection {
	return s.d.db.C(buildsC)
}

// GetAll, lists all builds belongs to an image by the matching owner and key
func (s *BuildStorage) GetAll(username string, imgKey string, opts domain.ListOptions) ([]domain.BuildRecord, domain.Pagination, error) {
	var build domain.Build
	err := s.col().Find(bson.M{"owner": username, "image_key": imgKey}).One(&build)
	if err != nil {
		return nil, domain.Pagination{}, toStorageErr(err)
	}

	nrecords := len(build.Records)
	skip, limit := opts.Cursor(nrecords)
	pagination := opts.Paginate(nrecords)

	return build.Records[skip:limit], pagination, nil
}

// GetLast, gets the latest record of a build by matching username and image key
func (s *BuildStorage) GetLast(username string, imgKey string) (domain.BuildRecord, error) {
	var build domain.Build
	err := s.col().Find(bson.M{"owner": username, "image_key": imgKey}).Select(bson.M{"records.$": 1}).One(&build)
	if err != nil || len(build.Records) == 0 {
		return domain.BuildRecord{}, domain.ErrNotFound
	}

	return build.Records[0], nil
}

// GetLastBy retrieves last build records for matching image keys
func (s *BuildStorage) GetLastBy(username string, imgKeys []string) (map[string]domain.BuildRecord, error) {
	var builds []domain.Build
	query := bson.M{"owner": username, "image_key": bson.M{"$in": imgKeys}, "records.0": bson.M{"$exists": true}}
	err := s.col().Find(query).Select(bson.M{"records": bson.M{"$slice": 1}}).All(&builds)
	if err != nil {
		return nil, toStorageErr(err)
	}

	records := make(map[string]domain.BuildRecord, len(builds))
	for _, b := range builds {
		if len(b.Records) == 0 {
			continue
		}
		records[b.ImageKey] = b.Records[0]
	}

	return records, nil
}

// List, lists builds of a user by matching username
func (s *BuildStorage) List(username string, opts domain.ListOptions) ([]domain.Build, domain.Pagination, error) {
	query := bson.M{"owner": username}

	nbuilds, err := s.col().Find(query).Count()
	if err != nil {
		return nil, domain.Pagination{}, toStorageErr(err)
	}

	skip, limit := opts.Cursor(nbuilds)
	pagination := opts.Paginate(nbuilds)

	var builds []domain.Build
	err = s.col().Find(bson.M{"owner": username, "$where": "this.records.length > 0"}).
		Sort("-last_build").
		Select(bson.M{"records": bson.M{"$slice": 1}}).
		Skip(skip).
		Limit(limit).
		All(&builds)

	return builds, pagination, toStorageErr(err)
}

// UpdateLast, updates the latest record of a build by matching username and image key
func (s *BuildStorage) UpdateLast(username string, imgKey string, record domain.BuildRecord) error {
	query := bson.M{"owner": username, "image_key": imgKey}
	update := bson.M{"$set": bson.M{"records.0": record}}
	err := s.col().Update(query, update)
	return toStorageErr(err)
}

// Put, puts a new build record as the latest record for a build by matching username and image key
func (s *BuildStorage) Put(username string, imgKey string, record domain.BuildRecord) error {
	query := bson.M{"owner": username, "image_key": imgKey}
	update := bson.M{"$push": bson.M{"records": bson.M{"$each": []domain.BuildRecord{record}, "$position": 0}}}
	err := s.col().Update(query, update)

	// If the build entry not found create one, that means this is the first
	// build for the image
	if err == mgo.ErrNotFound {
		now := time.Now()
		build := domain.Build{
			ImageKey:   imgKey,
			Owner:      username,
			LastRecord: now,
			Records:    []domain.BuildRecord{record},
		}
		err = s.col().Insert(build)
	}

	return toStorageErr(err)
}
