package mongo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Configuration is a structure of necessary information needed to run this
// service
type Configuration struct {
	// Conn is connection url for mongodb
	Conn string

	// Crt specifies the path to an x509 certificate file
	Crt string

	// Key specifies the path to the x509 key file
	Key string
}

// ConfigFromMap parses map input into Configuration
func ConfigFromMap(in map[string]interface{}) (*Configuration, error) {
	var config Configuration
	err := mapstructure.Decode(in, &config)
	return &config, err
}

type mongo struct {
	conn *mgo.Session
	db   *mgo.Database

	// TODO: Is it safe to keep these keys in memory?
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

const (
	day                       = time.Hour * 24
	authTokenValidDuration    = time.Minute * 15
	refreshTokenValidDuration = day * 5
)

var signingMethod = jwt.SigningMethodRS256

// New creates a mongodb backed authentication service
func New(ctx context.Context, timeout time.Duration, conf *Configuration) (auth.Service, error) {
	privateKeyBytes, err := ioutil.ReadFile(conf.Key)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := ioutil.ReadFile(conf.Crt)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	var sess *mgo.Session
	err = errs.RetryWithContext(ctx, timeout, time.Second*10, func() (err error) {
		logrus.Info("MongoDB auth trying to connect to the server...")
		sess, err = mgo.Dial(conf.Conn)
		return err
	})
	if err != nil {
		return nil, err
	}

	db := sess.DB("pullr")
	return &mongo{sess, db, signKey, verifyKey}, nil
}

func (a *mongo) Close() error {
	a.conn.Close()
	return nil
}

func (a *mongo) Validate(refreshToken, authToken string) (*auth.Secrets, string, error) {
	if refreshToken == "" || authToken == "" {
		return nil, "", auth.ErrUnauthenticated
	}

	claims := new(jwt.StandardClaims)
	authJwt, authParseErr := jwt.ParseWithClaims(authToken, claims, a.keyFunc)
	refreshJwt, refreshParseErr := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, a.keyFunc)
	if refreshParseErr != nil {
		return nil, "", auth.ErrUnauthenticated
	}

	var err error
	newAuthToken := authToken

	if !authJwt.Valid {
		ve, ok := authParseErr.(*jwt.ValidationError)
		expireErr := ok && ve.Errors&jwt.ValidationErrorExpired != 0
		if !expireErr {
			return nil, "", auth.ErrUnauthenticated
		}

		newAuthToken, err = a.updateAuthToken(refreshJwt, authJwt)
		if err != nil {
			return nil, "", err
		}
	}

	newRefreshToken, err := a.updateRefreshToken(refreshJwt)
	if err != nil {
		return nil, "", err
	}

	secrets := &auth.Secrets{
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
	}

	return secrets, claims.Subject, nil
}

func (a *mongo) Login(username, password string) (*auth.Secrets, error) {
	col := a.db.C("users")
	usr := new(domain.User)
	if err := col.Find(bson.M{"username": username}).One(usr); err != nil {
		if err == mgo.ErrNotFound {
			return nil, auth.ErrCredentials
		}

		return nil, err
	}

	if !usr.ComparePassword(password) {
		return nil, auth.ErrCredentials
	}

	authToken, err := a.createAuthToken(username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.createRefreshToken(username)
	if err != nil {
		return nil, err
	}

	secrets := &auth.Secrets{
		RefreshToken: refreshToken,
		AuthToken:    authToken,
	}

	return secrets, nil
}

func (a *mongo) Register(username, email, password string) error {
	users := a.db.C("users")
	numUsers, err := users.Find(bson.M{"username": username}).Count()
	if err != nil {
		return err
	}
	if numUsers > 0 {
		return auth.ErrUsernameTaken
	}

	numUsers, err = users.Find(bson.M{"email": email}).Count()
	if err != nil {
		return err
	}
	if numUsers > 0 {
		return auth.ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return users.Insert(bson.M{"username": username, "email": email, "password": hashedPassword})
}

func (a *mongo) ParseToken(token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, a.keyFunc)
}

func (a *mongo) SignToken(token *jwt.Token) (string, error) {
	return token.SignedString(a.signKey)
}

func (a *mongo) NewToken(claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(signingMethod, claims)
}

func (a *mongo) NewOAuthCbIdentifier(username, provider, redirectURI string) (auth.OAuthCbIdentifier, error) {
	ids := a.db.C("oauth_ids")
	newUUID := uuid.NewV1().String()
	cbIdentifier := auth.OAuthCbIdentifier{
		Provider:    provider,
		Username:    username,
		UUID:        newUUID,
		RedirectURI: redirectURI,
	}

	return cbIdentifier, ids.Insert(cbIdentifier)
}

func (a *mongo) OAuthCbIdentifier(uuid string) (*auth.OAuthCbIdentifier, error) {
	ids := a.db.C("oauth_ids")
	var cbIdentifier auth.OAuthCbIdentifier

	err := ids.Find(bson.M{"uuid": uuid}).One(&cbIdentifier)
	if err != nil {
		return nil, err
	}

	return &cbIdentifier, nil
}

func (a *mongo) RemoveOAuthCbIdentifier(uuid string) error {
	ids := a.db.C("oauth_ids")
	return ids.Remove(bson.M{"uuid": uuid})
}

func (a *mongo) keyFunc(token *jwt.Token) (interface{}, error) {
	return a.verifyKey, nil
}

func (a *mongo) createRefreshToken(username string) (string, error) {
	jti, err := randomString(32)
	if err != nil {
		return "", err
	}

	// FIXME: Auth - refresh token collisions
	if err := a.db.C("refresh_tokens").Insert(bson.M{"jti": jti}); err != nil {
		return "", err
	}

	tokenExp := time.Now().Add(refreshTokenValidDuration).Unix()

	claims := &jwt.StandardClaims{
		Id:        jti,
		Subject:   username,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return a.SignToken(token)
}

func (a *mongo) updateRefreshToken(oldToken *jwt.Token) (string, error) {
	oldClaims, ok := oldToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", auth.ErrInvalidToken
	}

	expire := time.Now().Add(refreshTokenValidDuration).Unix()
	newClaims := &jwt.StandardClaims{
		Id:        oldClaims.Id,
		Subject:   oldClaims.Subject,
		ExpiresAt: expire,
	}

	token := jwt.NewWithClaims(signingMethod, newClaims)
	return a.SignToken(token)
}

func (a *mongo) createAuthToken(username string) (string, error) {
	tokenExp := time.Now().Add(authTokenValidDuration).Unix()
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return a.SignToken(token)
}

func (a *mongo) updateAuthToken(refreshToken *jwt.Token, oldAuthToken *jwt.Token) (string, error) {
	refreshTokenClaims, ok := refreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", auth.ErrInvalidToken
	}

	if err := a.checkRefreshToken(refreshTokenClaims.Id); err != nil {
		return "", err
	}

	if !refreshToken.Valid {
		err := auth.ErrInvalidToken
		errs.Log(a.db.C("refresh_tokens").Remove(bson.M{"jti": refreshTokenClaims.Id}))
		return "", err
	}

	authTokenClaims, ok := oldAuthToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", auth.ErrInvalidToken
	}

	return a.createAuthToken(authTokenClaims.Subject)
}

func (a *mongo) checkRefreshToken(tokenID string) error {
	numToken, err := a.db.C("refresh_tokens").Find(bson.M{"jti": tokenID}).Count()
	if err != nil {
		return err
	}

	if numToken != 1 {
		return auth.ErrInvalidToken
	}

	return nil
}

func randomString(numBytes uint8) (string, error) {
	b := make([]byte, numBytes)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}

	str := base64.URLEncoding.EncodeToString(b)
	return str, nil
}
