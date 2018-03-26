package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AuthStorage stores and queries authentication data from mongodb
type AuthStorage struct {
	d *Driver
}

// internal struct for representing tokens in auth collection
type authToken struct {
	TokenType string `bson:"token_type"`
	TokenId   string `bson:"token_id"`
	Token     string `bson:"token"`
}

const tokenRefresh = "refresh"

func newAuthToken(tokenType, tokenId, token string) authToken {
	return authToken{tokenType, tokenId, token}
}

// internal struct for representing credentials in mongodb
type authCredentials struct {
	Username string      `bson:"username"`
	Password string      `bson:"password"`
	Email    string      `bson:"email"`
	Tokens   []authToken `bson:"tokens"`
}

func (s *AuthStorage) col() *mgo.Collection {
	return s.d.db.C(authC)
}

// GetPassword finds matching user's password
func (s *AuthStorage) GetPassword(username string) (string, error) {
	var credentials authCredentials
	err := s.col().Find(bson.M{"username": username}).Select(bson.M{"password": 1}).One(&credentials)
	return credentials.Password, toStorageErr(err)
}

// GetPasswordByEmail finds matching user's password by their email
func (s *AuthStorage) GetPasswordByEmail(email string) (string, error) {
	var credentials authCredentials
	err := s.col().Find(bson.M{"email": email}).Select(bson.M{"password": 1}).One(&credentials)
	return credentials.Password, toStorageErr(err)
}

// GetRefreshToken finds a refresh token by matching tokenID
func (s *AuthStorage) GetRefreshToken(tokenID string) (string, error) {
	query := bson.M{"tokens": bson.M{"$elemMatch": bson.M{"token_id": tokenID}}}
	proj := bson.M{"tokens.$": 1}

	var credentials authCredentials
	err := s.col().Find(query).Select(proj).One(&credentials)
	if err != nil {
		return "", toStorageErr(err)
	}

	return credentials.Tokens[0].Token, nil
}

// PutRefreshToken puts a refresh token into user credentials by matching username
func (s *AuthStorage) PutRefreshToken(username string, tokenID string, token string) error {
	query := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"tokens": newAuthToken(tokenRefresh, tokenID, token)}}
	err := s.col().Update(query, update)
	return toStorageErr(err)
}

// DeleteRefreshToken deletes a refresh token record by matching tokenID
func (s *AuthStorage) DeleteRefreshToken(tokenID string) error {
	update := bson.M{"$pull": bson.M{"tokens": bson.M{"token_id": tokenID, "token_type": tokenRefresh}}}
	err := s.col().Update(bson.M{}, update)
	return toStorageErr(err)
}

// PutCredentials puts a new credentials record into database
func (s *AuthStorage) PutCredentials(username string, email string, hashedPassword string) error {
	err := s.col().Insert(authCredentials{
		Tokens:   []authToken{},
		Password: hashedPassword,
		Username: username,
		Email:    email,
	})

	return toStorageErr(err)
}

// Delete deletes a credentials record by matching username
func (s *AuthStorage) Delete(username string) error {
	err := s.col().Remove(bson.M{"username": username})
	return toStorageErr(err)
}
