package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/api/auth"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// Api is v1 api wrapper
type Api struct {
	Group         *echo.Group
	Authenticator auth.Authenticator

	// Services
	buildsvc  *domain.BuildService
	authsvc   *domain.DefaultAuthService
	oauthsvc  *domain.OAuthService
	sourcesvc *domain.SourceService

	// Storages
	imageStorage domain.ImageStorage
	userStorage  domain.UserStorage
	buildStorage domain.BuildStorage
	oauthStorage domain.OAuthStorage
	authStorage  domain.AuthStorage
}

// NewApi add api v1 routes to the given routing group
func NewApi(config Config, authenticator auth.Authenticator, group *echo.Group) *Api {
	api := &Api{
		group,
		authenticator,
		config.BuildService,
		config.AuthService,
		config.OAuthService,
		config.SourceService,
		config.Storage.ImageStorage(),
		config.Storage.UserStorage(),
		config.Storage.BuildStorage(),
		config.Storage.OAuthStorage(),
		config.Storage.AuthStorage(),
	}

	// Authentication endpoints
	if config.HandleAuthentication {
		group.POST("/auth/login", api.AuthLogin)
		group.POST("/auth/register", api.AuthRegister)
	}

	// Restricted group routes are protected and needs authentication
	restricted := group.Group("", authenticator.Middleware())

	// OAuth endpoints
	if config.HandleOAuth {
		restricted.GET("/oauth/:provider/login_url", authenticator.Wrap(api.OAuthLogin))
		group.GET("/oauth/:provider/cb/:username", api.OAuthCallback)
	}

	// User endpoints
	restricted.GET("/user/profile", authenticator.Wrap(api.UserProfile))
	restricted.POST("/user/profile", authenticator.Wrap(api.UserProfileUpdate))

	// Image endpoints
	restricted.GET("/images", authenticator.Wrap(api.ImageList))
	restricted.POST("/images", authenticator.Wrap(api.ImageCreate))
	restricted.GET("/images/:key", authenticator.Wrap(api.ImageGet))
	restricted.POST("/images/:key", authenticator.Wrap(api.ImageUpdate))
	restricted.DELETE("/images/:key", authenticator.Wrap(api.ImageDelete))

	// Build endpoints
	restricted.GET("/builds", authenticator.Wrap(api.BuildList))
	restricted.GET("/builds/:key", authenticator.Wrap(api.BuildHistory))

	// SourceClient endpoints
	group.POST("/source/:provider/:username/webhook", api.SourceWebhook)
	restricted.GET("/source/:provider/orgs", authenticator.Wrap(api.SourceOrganisations))
	restricted.GET("/source/:provider/repos", authenticator.Wrap(api.SourceRepositories))

	return api
}
