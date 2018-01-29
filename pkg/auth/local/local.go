package local

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mobingilabs/pullr/pkg/auth"
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
	AuthTokenValidDuration    = time.Minute * 15
	RefreshTokenValidDuration = day * 5
)

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
func (a *Authenticator) Validate(csrf, refreshToken, authToken string) (*auth.Secrets, string, error) {
	if csrf == "" || refreshToken == "" || authToken == "" {
		return nil, "", auth.ErrUnauthenticated
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return a.verifyKey, nil
	}

	authJwt, authParseErr := jwt.ParseWithClaims(authToken, &auth.TokenClaims{}, keyFunc)
	if authParseErr != nil {
		return nil, "", authParseErr
	}

	refreshJwt, refreshParseErr := jwt.ParseWithClaims(refreshToken, &auth.TokenClaims{}, keyFunc)
	if refreshParseErr != nil {
		return nil, "", refreshParseErr
	}

	authClaims, ok := authJwt.Claims.(*auth.TokenClaims)
	if !ok {
		return nil, "", auth.ErrInvalidToken
	}

	if csrf != authClaims.Csrf {
		return nil, "", auth.ErrUnauthenticated
	}

	var err error
	newAuthToken := authToken
	newCsrf := csrf
	newRefreshToken := ""

	if !authJwt.Valid {
		if !isTokenExpireErr(authParseErr) {
			return nil, "", auth.ErrUnauthenticated
		}

		newAuthToken, newCsrf, err = a.updateAuthToken(refreshJwt, authJwt)
		if err != nil {
			return nil, "", err
		}
	}

	newRefreshToken, err = a.updateRefreshToken(refreshJwt, newCsrf)
	if err != nil {
		return nil, "", err
	}

	secrets := &auth.Secrets{
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
		Csrf:         newCsrf,
	}

	return secrets, authClaims.Subject, nil
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

	csrf, err := randomString(32)
	if err != nil {
		return nil, err
	}

	authToken, err := a.createAuthToken(username, csrf)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.createRefreshToken(username, csrf)
	if err != nil {
		return nil, err
	}

	secrets := &auth.Secrets{
		RefreshToken: refreshToken,
		AuthToken:    authToken,
		Csrf:         csrf,
	}

	return secrets, nil
}

func (a *Authenticator) Register(username, password string) error {
	// TODO: Look example to implement register
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

func (a *Authenticator) Sign(token *jwt.Token) (string, error) {
	return token.SignedString(a.signKey)
}

func (a *Authenticator) createRefreshToken(username, csrf string) (string, error) {
	jti, err := randomString(32)
	if err != nil {
		return "", err
	}

	// FIXME: Auth - refresh token collisions
	if err := a.db.C("refresh_tokens").Insert(bson.M{"jti": jti}); err != nil {
		return "", err
	}

	tokenExp := time.Now().Add(RefreshTokenValidDuration).Unix()

	claims := auth.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   username,
			ExpiresAt: tokenExp,
		},
		Csrf: csrf,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return a.Sign(token)
}

func (a *Authenticator) updateRefreshToken(oldToken *jwt.Token, csrf string) (string, error) {
	oldClaims, ok := oldToken.Claims.(*auth.TokenClaims)
	if !ok {
		return "", auth.ErrInvalidToken
	}

	expire := time.Now().Add(RefreshTokenValidDuration).Unix()
	newClaims := auth.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        oldClaims.Id,
			Subject:   oldClaims.Subject,
			ExpiresAt: expire,
		},
		Csrf: csrf,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, newClaims)
	return a.Sign(token)
}

func (a *Authenticator) createAuthToken(username string, csrf string) (string, error) {
	tokenExp := time.Now().Add(AuthTokenValidDuration).Unix()
	claims := auth.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: tokenExp,
		},
		Csrf: csrf,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return a.Sign(token)
}

func (a *Authenticator) updateAuthToken(refreshToken *jwt.Token, oldAuthToken *jwt.Token) (newToken, csrf string, err error) {
	refreshTokenClaims, ok := refreshToken.Claims.(*auth.TokenClaims)
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

	authTokenClaims, ok := oldAuthToken.Claims.(*auth.TokenClaims)
	if !ok {
		err = auth.ErrInvalidToken
		return
	}

	csrf, err = randomString(32)
	if err != nil {
		return
	}

	newToken, err = a.createAuthToken(authTokenClaims.Subject, authTokenClaims.Csrf)
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

func isTokenExpireErr(err error) bool {
	ve, ok := err.(*jwt.ValidationError)
	return ok && ve.Errors&jwt.ValidationErrorExpired != 0
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
