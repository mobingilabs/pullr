package local

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Authenticator implements a local database based authentication over mongodb
type Authenticator struct {
	conn *mgo.Session
	db   *mgo.Database

	// TODO: Is it safe to keep these keys in memory?
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

// User credentials for auth
type User struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
}

const (
	day                       = time.Hour * 24
	authTokenValidDuration    = time.Minute * 15
	refreshTokenValidDuration = day * 5
)

var signingMethod = jwt.SigningMethodRS256

// New creates a new Authenticator instance with given mongodb connection,
// mongodb connection should be unique to the authenticator, make sure the
// connection is not shared with other services.
func New(conn *mgo.Session, privKeyPath string, pubKeyPath string) (*Authenticator, error) {
	privateKeyBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	db := conn.DB("pullr")
	return &Authenticator{conn, db, signKey, verifyKey}, nil
}

// Close closes the mongodb connection
func (a *Authenticator) Close() {
	a.conn.Close()
}

// Validate reports back the extracted username from the JWT token if the token
// is valid.
func (a *Authenticator) Validate(refreshToken, authToken string) (*auth.Secrets, string, error) {
	if refreshToken == "" || authToken == "" {
		return nil, "", auth.ErrUnauthenticated
	}

	claims := new(jwt.StandardClaims)
	authJwt, authParseErr := jwt.ParseWithClaims(authToken, claims, a.keyFunc)
	refreshJwt, refreshParseErr := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, a.keyFunc)
	if refreshParseErr != nil {
		return nil, "", refreshParseErr
	}

	var err error
	newAuthToken := authToken
	newRefreshToken := ""

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

	newRefreshToken, err = a.updateRefreshToken(refreshJwt)
	if err != nil {
		return nil, "", err
	}

	secrets := &auth.Secrets{
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
	}

	return secrets, claims.Subject, nil
}

// Login will generate a new JWT token for the given user
func (a *Authenticator) Login(username, password string) (*auth.Secrets, error) {
	col := a.db.C("users")
	var user User
	err := col.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, auth.ErrCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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

func (a *Authenticator) Register(username, password string) error {
	users := a.db.C("users")
	numUsers, err := users.Find(bson.M{"username": username}).Count()
	if err != nil {
		return err
	}

	if numUsers > 0 {
		return auth.ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err := users.Insert(bson.M{"username": username, "password": hashedPassword}); err != nil {
		return err
	}

	return nil
}

func (a *Authenticator) ParseToken(token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, a.keyFunc)
}

func (a *Authenticator) SignToken(token *jwt.Token) (string, error) {
	return token.SignedString(a.signKey)
}

func (a *Authenticator) NewToken(claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(signingMethod, claims)
}

// NewOAuthCbIdentifier generates an identifier for the given user to use with
// oauth providers login mechanism
func (a *Authenticator) NewOAuthCbIdentifier(username, provider, redirectUri string) (auth.OAuthCbIdentifier, error) {
	ids := a.db.C("oauth_ids")
	newId := uuid.NewV1().String()
	cbIdentifier := auth.OAuthCbIdentifier{
		Provider:    provider,
		Username:    username,
		Uuid:        newId,
		RedirectUri: redirectUri,
	}

	return cbIdentifier, ids.Insert(cbIdentifier)
}

// OAuthUserFromIdentifier reports back the user identity from given oauth
// identifier
func (a *Authenticator) OAuthCbIdentifier(uuid string) (*auth.OAuthCbIdentifier, error) {
	ids := a.db.C("oauth_ids")
	var cbIdentifier auth.OAuthCbIdentifier

	err := ids.Find(bson.M{"uuid": uuid}).One(&cbIdentifier)
	if err != nil {
		return nil, err
	}

	return &cbIdentifier, nil
}

func (a *Authenticator) RemoveOAuthCbIdentifier(uuid string) error {
	ids := a.db.C("oauth_ids")
	return ids.Remove(bson.M{"uuid": uuid})
}

func (a *Authenticator) keyFunc(token *jwt.Token) (interface{}, error) {
	return a.verifyKey, nil
}

func (a *Authenticator) createRefreshToken(username string) (string, error) {
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

func (a *Authenticator) updateRefreshToken(oldToken *jwt.Token) (string, error) {
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

func (a *Authenticator) createAuthToken(username string) (string, error) {
	tokenExp := time.Now().Add(authTokenValidDuration).Unix()
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return a.SignToken(token)
}

func (a *Authenticator) updateAuthToken(refreshToken *jwt.Token, oldAuthToken *jwt.Token) (newToken string, err error) {
	refreshTokenClaims, ok := refreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		err = auth.ErrInvalidToken
		return
	}

	if err = a.checkRefreshToken(refreshTokenClaims.Id); err != nil {
		return
	}

	if !refreshToken.Valid {
		err = auth.ErrInvalidToken
		a.db.C("refresh_tokens").Remove(bson.M{"jti": refreshTokenClaims.Id})
		return
	}

	authTokenClaims, ok := oldAuthToken.Claims.(*jwt.StandardClaims)
	if !ok {
		err = auth.ErrInvalidToken
		return
	}

	newToken, err = a.createAuthToken(authTokenClaims.Subject)
	return
}

func (a *Authenticator) checkRefreshToken(tokenId string) error {
	numToken, err := a.db.C("refresh_tokens").Find(bson.M{"jti": tokenId}).Count()
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
