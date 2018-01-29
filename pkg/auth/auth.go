package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrCredentials     = errors.New("credentials are not met")
	ErrUsernameTaken   = errors.New("username has already taken")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
)

type TokenClaims struct {
	jwt.StandardClaims
	Csrf string `json:"csrf"`
}

// Secrets represents all the tokens required for identifying the subject
type Secrets struct {
	RefreshToken string
	AuthToken    string
	Csrf         string
}

// Authenticator handles token based authentication.
//
// By default Pullr assumes tokens are JWTs and will be sent to the client in
// the response body, please keep that in mind and never expose any user secrets
// with the token.
type Authenticator interface {
	// Validate checks if the given tokens are valid and then updates the tokens
	// for further requests and also returns token's subject
	Validate(csrf, refreshToken, authToken string) (*Secrets, string, error)

	// Login will generate a token if the given credentials are correct for the
	// user.
	Login(username, password string) (*Secrets, error)

	// Register will create user record
	Register(username, password string) error

	// Sign reports signed token as string
	Sign(token *jwt.Token) (string, error)
}

// Token describes the information kept in the generated token
type Token struct {
	// Valid is true if the parsed token is validated
	Valid bool
	// Username is empty if the token is not valid
	Username string
}
