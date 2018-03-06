package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
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

type mongo struct {
	session *mgo.Session
	db      *mgo.Database
}

// ConfigFromMap parses a map into Config
func ConfigFromMap(in map[string]interface{}) (*Config, error) {
	var config Config
	err := mapstructure.Decode(in, &config)
	return &config, err
}

func New(ctx context.Context, timeout time.Duration, conf *Config) (*mongo, error) {
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

func (m *mongo) Close() error {
	m.session.Close()
	return nil
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
	if opts == nil {
		return "$natural"
	}

	s := opts.SortBy
	if s == "" {
		return "$natural"
	}

	dirSign := "-"
	if opts.SortDirection == storage.Asc {
		dirSign = "+"
	}

	return fmt.Sprintf("%s%s", dirSign, s)
}

func mergeBson(src bson.M, others ...bson.M) bson.M {
	clone := make(bson.M, len(src))
	for k := range src {
		clone[k] = src[k]
	}

	for _, other := range others {
		for k := range other {
			clone[k] = other[k]
		}
	}

	return clone
}
