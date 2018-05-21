package v1

import "github.com/mobingilabs/pullr/pkg/domain"

// Config can be used for configuring api server's behaviour
type Config struct {
	// HandleAuthentication, if true /auth/* endpoints will be activated
	// Default: true
	HandleAuthentication bool

	// HandleOAuth, if true pullr will be responsible for linking oauth
	// accounts.
	// Default: true
	HandleOAuth bool

	Storage       domain.StorageDriver
	BuildService  *domain.BuildService
	AuthService   *domain.DefaultAuthService
	OAuthService  *domain.OAuthService
	SourceService *domain.SourceService
}

// NewConfig creates an api configuration object with defaults
func NewConfig() Config {
	return Config{HandleAuthentication: true, HandleOAuth: true}
}
