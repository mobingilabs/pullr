package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	day                       = time.Hour * 24
	authTokenValidDuration    = time.Minute * 15
	refreshTokenValidDuration = day * 5
)

var signingMethod = jwt.SigningMethodRS256
var noSecrets AuthSecrets

// AuthStorage stores authentication data such as tokens or credentials
type AuthStorage interface {
	GetPassword(username string) (string, error)
	GetPasswordByEmail(email string) (string, error)
	GetRefreshToken(tokenID string) (string, error)
	PutRefreshToken(username string, tokenID string, token string) error
	DeleteRefreshToken(tokenID string) error
	PutCredentials(username string, email string, hashedPassword string) error
	Delete(username string) error
}

// AuthSecrets is combined structure of auth token and refresh token
type AuthSecrets struct {
	Username     string `json:"username"`
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

// DefaultAuthService authenticates and grants users
type DefaultAuthService struct {
	storage   AuthStorage
	users     UserStorage
	log       Logger
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

// NewAuthService creates a default authentication service. Auth service is responsible for
// registering, logging in users and providing secrets for them. Auth service itself
// is not responsible for storing the user information, it keeps only credentials
// data for the users. Persistence backend can be configured by the storage parameter.
func NewAuthService(storage AuthStorage, users UserStorage, logger Logger, conf AuthConfig) (*DefaultAuthService, error) {
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

	service := &DefaultAuthService{
		storage:   storage,
		log:       logger,
		users:     users,
		signKey:   signKey,
		verifyKey: verifyKey,
	}

	return service, nil
}

// Grant, validates the secrets given by the previously authenticated user and grants
// access for the requested resources with updated secrets
func (s *DefaultAuthService) Grant(refreshToken, authToken string) (AuthSecrets, error) {
	if refreshToken == "" || authToken == "" {
		return noSecrets, ErrAuthUnauthorized
	}

	claims := new(jwt.StandardClaims)
	authJwt, authParseErr := jwt.ParseWithClaims(authToken, claims, s.keyFunc)
	refreshJwt, refreshParseErr := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, s.keyFunc)
	if refreshParseErr != nil {
		return noSecrets, ErrAuthUnauthorized
	}

	var err error
	newAuthToken := authToken

	if authParseErr != nil || !authJwt.Valid {
		ve, ok := authParseErr.(*jwt.ValidationError)
		expireErr := ok && ve.Errors&jwt.ValidationErrorExpired != 0
		if !expireErr {
			return noSecrets, ErrAuthUnauthorized
		}

		newAuthToken, err = s.updateAuthToken(refreshJwt, authJwt)
		if err != nil {
			return noSecrets, err
		}
	}

	newRefreshToken, err := s.updateRefreshToken(refreshJwt)
	if err != nil {
		return noSecrets, err
	}
	secrets := AuthSecrets{
		Username:     claims.Subject,
		AuthToken:    newAuthToken,
		RefreshToken: newRefreshToken,
	}

	return secrets, nil
}

// Login authenticates a user if given username and password matches the authentication
// records
func (s *DefaultAuthService) Login(username, password string) (AuthSecrets, error) {
	pass, err := s.storage.GetPassword(username)
	if err != nil {
		return AuthSecrets{}, ErrAuthBadCredentials
	}

	matchErr := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	if matchErr != nil {
		return AuthSecrets{}, ErrAuthBadCredentials
	}

	authToken, err := s.createAuthToken(username)
	if err != nil {
		return AuthSecrets{}, err
	}

	refreshToken, err := s.createRefreshToken(username)
	if err != nil {
		return AuthSecrets{}, err
	}

	secrets := AuthSecrets{
		Username:     username,
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}

	return secrets, nil
}

// Register, registers a new user with given credentials. It is not meant
// to save any profile related information. Only credentials are saved.
func (s *DefaultAuthService) Register(username, email, password string) error {
	_, err := s.storage.GetPassword(username)

	// If username is already exists do not continue
	if err == nil {
		return ErrUserUsernameExist
	} else if err != ErrNotFound {
		return err
	}

	_, err = s.storage.GetPasswordByEmail(email)

	// If email is already exists do not continue
	if err == nil {
		return ErrUserEmailExist
	} else if err != ErrNotFound {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.storage.PutCredentials(username, email, string(hashedPassword))
}

func (s *DefaultAuthService) signToken(token *jwt.Token) (string, error) {
	return token.SignedString(s.signKey)
}

func (s *DefaultAuthService) keyFunc(token *jwt.Token) (interface{}, error) {
	return s.verifyKey, nil
}

func (s *DefaultAuthService) createRefreshToken(username string) (string, error) {
	jti, err := randomString(32)
	if err != nil {
		return "", err
	}

	tokenExp := time.Now().Add(refreshTokenValidDuration).Unix()

	claims := &jwt.StandardClaims{
		Id:        jti,
		Subject:   username,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	signedToken, err := s.signToken(token)
	if err != nil {
		return "", err
	}

	err = s.storage.PutRefreshToken(username, jti, signedToken)
	return signedToken, err
}

func (s *DefaultAuthService) updateRefreshToken(oldToken *jwt.Token) (string, error) {
	oldClaims, ok := oldToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", ErrAuthBadToken
	}

	expire := time.Now().Add(refreshTokenValidDuration).Unix()
	newClaims := &jwt.StandardClaims{
		Id:        oldClaims.Id,
		Subject:   oldClaims.Subject,
		ExpiresAt: expire,
	}

	token := jwt.NewWithClaims(signingMethod, newClaims)
	return s.signToken(token)
}

func (s *DefaultAuthService) createAuthToken(username string) (string, error) {
	tokenExp := time.Now().Add(authTokenValidDuration).Unix()
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return s.signToken(token)
}

func (s *DefaultAuthService) updateAuthToken(refreshToken *jwt.Token, oldAuthToken *jwt.Token) (string, error) {
	refreshTokenClaims, ok := refreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", ErrAuthBadToken
	}

	if _, err := s.storage.GetRefreshToken(refreshTokenClaims.Id); err != nil {
		return "", ErrAuthBadToken
	}

	authTokenClaims, ok := oldAuthToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", ErrAuthBadToken
	}

	return s.createAuthToken(authTokenClaims.Subject)
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
