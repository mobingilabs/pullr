package mongodb

import (
	"github.com/mobingilabs/pullr/pkg/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OAuthStorage stores and queries oauth token data from mongodb
type OAuthStorage struct {
	d *Driver
}

// oauthRecord is the internal representation of mongodb document
type oauthRecord struct {
	Username string              `bson:"username"`
	Tokens   []domain.OAuthToken `bson:"tokens"`
	Secrets  []oauthSecret       `bson:"secrets"`
}

type oauthSecret struct {
	Secret   string `bson:"secret"`
	Redirect string `bson:"redirect"`
}

func (s *OAuthStorage) col() *mgo.Collection {
	return s.d.db.C(oauthC)
}

// PutSecret, puts a new oauth login secret into the database
func (s *OAuthStorage) PutSecret(username, secret, cburi string) error {
	query := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"secrets": oauthSecret{secret, cburi}}}
	_, err := s.col().Upsert(query, update)
	return toStorageErr(err)
}

// PopRedirectURL, finds a secret record by matching username and secret and,
// reports back the associated redirect url. If it is found secret record
// will be deleted.
func (s *OAuthStorage) PopRedirectURL(username, secret string) (string, error) {
	query := bson.M{"username": username, "secrets": bson.M{"$elemMatch": bson.M{"secret": secret}}}
	proj := bson.M{"secrets.$": 1}

	var record oauthRecord
	err := s.col().Find(query).Select(proj).One(&record)
	if err != nil {
		return "", err
	}

	return record.Secrets[0].Redirect, nil
}

// GetTokens, finds oauth tokens associated with a user by matching username
func (s *OAuthStorage) GetTokens(username string) (map[string]domain.OAuthToken, error) {
	query := bson.M{"username": username}

	var record oauthRecord
	err := s.col().Find(query).One(&record)
	if err != nil {
		return nil, err
	}
	if record.Tokens == nil {
		return make(map[string]domain.OAuthToken), nil
	}

	tokens := make(map[string]domain.OAuthToken, len(record.Tokens))
	for _, t := range record.Tokens {
		tokens[t.Provider] = t
	}

	return tokens, nil
}

// PutToken puts a new oauth token into user record by matching username
func (s *OAuthStorage) PutToken(username string, identity string, provider string, token string) error {
	tokenRecord := domain.OAuthToken{
		Provider: provider,
		Identity: identity,
		Token:    token,
	}

	query := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"tokens": tokenRecord}}
	_, err := s.col().Upsert(query, update)
	return toStorageErr(err)
}

// RemoveToken remove an oauth token from user's record by matching its provider
func (s *OAuthStorage) RemoveToken(username string, provider string) error {
	query := bson.M{"username": username}
	update := bson.M{"$pull": bson.M{"tokens": bson.M{"provider": provider}}}
	err := s.col().Update(query, update)
	return toStorageErr(err)
}
