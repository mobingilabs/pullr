package mongo

import (
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) UpdateStatus(status domain.Status) error {
	return m.db.C("history").Insert(status)
}

func (m *mongo) Status(username, kind, id string) (*domain.Status, error) {
	var status domain.Status
	err := m.db.C("history").Find(bson.M{"kind": kind, "id": id, "account": username}).Sort("-time").One(&status)
	return &status, err
}

func (m *mongo) Statuses(username string, kind string, listOpts *storage.ListOptions) ([]domain.Status, error) {
	col := m.db.C("history")
	if listOpts == nil {
		listOpts = storage.NewListOptions()
	}

	query := bson.M{"kind": kind, "account": username}
	skip := listOpts.Page * listOpts.PerPage
	limit := listOpts.PerPage

	pipe := []bson.M{
		{
			"$match": query,
		},
		{"$sort": bson.M{"time": -1}},
		{
			"$group": bson.M{
				"_id": "$id",
				"doc": bson.M{"$first": "$$ROOT"},
			},
		},
		{"$limit": limit},
		{"$skip": skip},
		{"$sort": bson.M{"time": -1}},
	}

	var results []struct {
		Doc domain.Status `bson:"doc"`
	}

	if err := col.Pipe(pipe).All(&results); err != nil {
		return nil, errors.WithMessage(err, "mongo storage failed to get statuses")
	}

	statuses := make([]domain.Status, len(results))
	for i, r := range results {
		statuses[i] = r.Doc
	}

	return statuses, nil
}

func (m *mongo) StatusesByResources(username string, kind string, ids []string) ([]domain.Status, error) {
	col := m.db.C("history")
	query := bson.M{"kind": kind, "account": username, "id": bson.M{"$in": ids}}
	pipe := []bson.M{
		{
			"$match": query,
		},
		{"$sort": bson.M{"time": -1}},
		{
			"$group": bson.M{
				"_id": "$id",
				"doc": bson.M{"$first": "$$ROOT"},
			},
		},
	}

	var results []struct {
		Doc domain.Status `bson:"doc"`
	}

	if err := col.Pipe(pipe).All(&results); err != nil {
		return nil, errors.WithMessage(err, "storage mongodb failed to get statuses for resource ids")
	}

	statuses := make([]domain.Status, len(results))
	for i, r := range results {
		statuses[i] = r.Doc
	}

	return statuses, nil
}

func (m *mongo) StatusesByCause(username, kind, cause string, listOpts *storage.ListOptions) ([]domain.Status, error) {
	col := m.db.C("history")
	if listOpts == nil {
		listOpts = storage.NewListOptions()
	}

	skip := listOpts.Page * listOpts.PerPage
	limit := listOpts.PerPage

	query := bson.M{"kind": kind, "cause": cause, "account": username}
	pipe := []bson.M{
		{
			"$match": query,
		},
		{"$sort": bson.M{"time": -1}},
		{
			"$group": bson.M{
				"_id": "$id",
				"doc": bson.M{"$first": "$$ROOT"},
			},
		},
		{"$limit": limit},
		{"$skip": skip},
	}

	var results []struct {
		Doc domain.Status `bson:"doc"`
	}

	if err := col.Pipe(pipe).All(&results); err != nil {
		return nil, err
	}

	statuses := make([]domain.Status, len(results))
	for i, r := range results {
		statuses[i] = r.Doc
	}
	return statuses, nil
}

func (m *mongo) History(username string, kind string, id string, listOpts *storage.ListOptions) ([]domain.Status, error) {
	var statuses []domain.Status
	if listOpts == nil {
		listOpts = storage.NewListOptions()
	}

	col := m.db.C("history")
	query := bson.M{"kind": kind, "id": id, "account": username}
	skip := listOpts.Page * listOpts.PerPage
	limit := listOpts.PerPage
	err := col.Find(query).Skip(skip).Limit(limit).Sort("-time").All(&statuses)
	return statuses, err
}
