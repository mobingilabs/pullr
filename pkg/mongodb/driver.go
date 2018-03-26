package mongodb

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/run"
	"gopkg.in/mgo.v2"
)

// collection names for stores
const (
	usersC  = "users"
	imagesC = "images"
	buildsC = "builds"
	authC   = "user_creds"
	oauthC  = "oauth"
)

// Config is a structure of necessary information needed to run this
// service
type Config struct {
	Conn string
}

// Driver is mongodb storage driver
type Driver struct {
	session *mgo.Session
	db      *mgo.Database
	logger  domain.Logger
}

// ConfigFromMap parses a map into Config
func ConfigFromMap(in map[string]string) (*Config, error) {
	var config Config
	err := mapstructure.Decode(in, &config)
	return &config, err
}

// Dial, creates a mongodb StorageDriver by dialing the mongodb host
func Dial(ctx context.Context, logger domain.Logger, conf *Config) (*Driver, error) {
	var sess *mgo.Session
	err := run.RetryWithContext(ctx, time.Second*10, func() (err error) {
		logger.Info("MongoDB storage trying to connect to the server...")
		sess, err = mgo.Dial(conf.Conn)
		return err
	})
	if err != nil {
		return nil, err
	}

	mongodb := Driver{
		session: sess,
		db:      sess.DB("pullr"),
		logger:  logger,
	}

	return &mongodb, nil
}

// Close closes the mongodb connection
func (d *Driver) Close() error {
	d.session.Close()
	return nil
}

// AuthStorage creates a mongodb baked AuthStorage
func (d *Driver) AuthStorage() domain.AuthStorage {
	return &AuthStorage{d}
}

// OAuthStorage creates a mongodb baked OAuthStorage
func (d *Driver) OAuthStorage() domain.OAuthStorage {
	return &OAuthStorage{d}
}

// UserStorage creates a mongodb baked UserStorage
func (d *Driver) UserStorage() domain.UserStorage {
	return &UserStorage{d}
}

// ImageStorage creates a mongodb baked ImageStorage
func (d *Driver) ImageStorage() domain.ImageStorage {
	return &ImageStorage{d}
}

// BuildStorage creates a mongodb baked BuildStorage
func (d *Driver) BuildStorage() domain.BuildStorage {
	return &BuildStorage{d}
}

func toStorageErr(err error) error {
	switch err {
	case nil:
		return nil
	case mgo.ErrNotFound:
		return domain.ErrNotFound
	default:
		return domain.ErrStorageDriver
	}
}
