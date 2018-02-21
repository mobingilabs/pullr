package mongo

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/storage"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersC  = "users"
	imagesC = "images"
)

// Config is a structure of necessary information needed to run this
// service
type Config struct {
	Conn string
}

// ConfigFromMap parses a map into Config
func ConfigFromMap(in map[string]interface{}) (*Config, error) {
	var config Config
	err := mapstructure.Decode(in, &config)
	return &config, err
}

type mongo struct {
	session *mgo.Session
	db      *mgo.Database
}

// New creates a mongodb backed storage service
func New(ctx context.Context, timeout time.Duration, conf *Config) (storage.Service, error) {
	var sess *mgo.Session
	err := errs.RetryWithContext(ctx, timeout, time.Second*10, func() (err error) {
		log.Info("MongoDB storage trying to connect to the server...")
		sess, err = mgo.Dial(conf.Conn)
		return err
	})
	if err != nil {
		return nil, err
	}

	mongodb := mongo{
		session: sess,
		db:      sess.DB("pullr"),
	}

	return &mongodb, nil
}

func (s *mongo) Close() error {
	s.session.Close()
	return nil
}

func (s *mongo) FindImageByKey(key string) (domain.Image, error) {
	col := s.db.C(imagesC)

	var image domain.Image
	query := bson.M{"key": key}
	err := col.Find(query).One(&image)
	return image, toStorageErr(err)
}

func (s *mongo) FindAllImages(username string, opts *storage.ListOptions) ([]domain.Image, storage.Pagination, error) {
	col := s.db.C(imagesC)

	var pagination storage.Pagination

	query := bson.M{"owner": username}
	count, err := col.Find(query).Count()
	if err != nil && err != mgo.ErrNotFound {
		return nil, pagination, err
	}

	page := opts.GetPage()
	perPage := opts.GetPerPage()
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
	err = col.Find(query).Sort(optsToMongoSort(opts)).Limit(perPage).Skip(perPage * page).All(&images)

	return images, pagination, toStorageErr(err)
}

func (s *mongo) FindAllImagesSince(username string, since time.Time) ([]domain.Image, error) {
	col := s.db.C(imagesC)

	var images []domain.Image
	query := bson.M{
		"owner": username,
		"$or": []bson.M{
			{"updated_at": bson.M{"$gt": since}},
			{"created_at": bson.M{"$gt": since}},
		},
	}
	err := col.Find(query).Sort("name").All(&images)

	return images, err
}

func (s *mongo) FindUser(username string) (domain.User, error) {
	col := s.db.C(usersC)

	var user domain.User
	err := col.Find(bson.M{"username": username}).One(&user)
	return user, toStorageErr(err)
}

func (s *mongo) PutUserToken(username, provider, token string) error {
	usr, err := s.FindUser(username)
	if err != nil {
		return err
	}

	if usr.Tokens == nil {
		usr.Tokens = make(map[string]domain.UserToken)
	}

	usr.PutToken(provider, username, token)
	return s.UpdateUser(username, usr)
}

func (s *mongo) CreateImage(image domain.Image) (string, error) {
	image.Key = domain.ImageKey(image.Repository)
	err := s.db.C(imagesC).Insert(image)
	return image.Key, toStorageErr(err)
}

func (s *mongo) UpdateImage(oldKey string, image domain.Image) (string, error) {
	newKey := domain.ImageKey(image.Repository)
	if newKey == "" {
		newKey = oldKey
	}

	image.Key = newKey
	err := s.db.C(imagesC).Update(bson.M{"key": oldKey}, bson.M{"$set": image})
	return newKey, toStorageErr(err)
}

func (s *mongo) DeleteImage(imageKey string) error {
	err := s.db.C(imagesC).Remove(bson.M{"key": imageKey})
	return toStorageErr(err)
}

func (s *mongo) CreateUser(user domain.User) error {
	err := s.db.C(usersC).Insert(user)
	return toStorageErr(err)
}

func (s *mongo) UpdateUser(username string, user domain.User) error {
	err := s.db.C(usersC).Update(bson.M{"username": username}, user)
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

func optsToMongoSort(opts *storage.ListOptions) string {
	s := opts.GetSortBy()
	if s == "" {
		return "$natural"
	}

	dirSign := "-"
	if opts.GetSortDirection() == storage.Asc {
		dirSign = "+"
	}

	return fmt.Sprintf("%s%s", dirSign, s)
}
