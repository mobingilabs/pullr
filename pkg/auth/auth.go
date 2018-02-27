package auth

import (
	"errors"
	"io"

	"github.com/dgrijalva/jwt-go"
)

// Authentication errors
var (
	ErrCredentials     = errors.New("credentials are not met")
	ErrUsernameTaken   = errors.New("username has already taken")
	ErrEmailTaken      = errors.New("this email has been already registered")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
)

// Secrets represents all the tokens required for identifying the subject
type Secrets struct {
	RefreshToken string
	AuthToken    string
}

// Service handles token based authentication.
//
// By default Pullr assumes tokens are JWTs and will be sent to the client in
// the response, please keep that in mind and never expose any user secrets
// with the token.
type Service interface {
	io.Closer

	// Validate checks if the given tokens are valid and then updates the tokens
	// for further requests and also returns token's subject
	Validate(refreshToken, authToken string) (*Secrets, string, error)

	// Login will generate a token if the given credentials are correct for the
	// user.
	Login(username, password string) (*Secrets, error)

	// Register will create user record
	Register(username, email, password string) error

	// ParseToken parses signed token
	ParseToken(token string, claims jwt.Claims) (*jwt.Token, error)

	// SignToken signs a given token
	SignToken(token *jwt.Token) (string, error)

	// NewToken creates a jwt token with given claims
	NewToken(claims jwt.Claims) *jwt.Token

	// NewOAuthCbIdentifier generates an identifier for the given user to use with
	// oauth providers login mechanism
	NewOAuthCbIdentifier(username, provider, redirectURI string) (OAuthCbIdentifier, error)

	// OAuthUserFromIdentifier reports back the user identity from given oauth
	// identifier
	OAuthCbIdentifier(uuid string) (*OAuthCbIdentifier, error)

	// RemoveOAuthCbIdentitifer remove identifier record
	RemoveOAuthCbIdentifier(uuid string) error
}

// OAuthCbIdentifier is used for identifying incoming oauth provider callbacks
type OAuthCbIdentifier struct {
	Username    string `bson:"username"`
	Provider    string `bson:"provider"`
	UUID        string `bson:"uuid"`
	RedirectURI string `bson:"redirect_uri"`
}
