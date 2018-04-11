package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/api/auth"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// Api is v1 api wrapper
type Api struct {
	Group *echo.Group

	// Services
	buildsvc  *domain.BuildService
	authsvc   *domain.AuthService
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
func NewApi(storage domain.StorageDriver, buildsvc *domain.BuildService, authsvc *domain.AuthService, oauthsvc *domain.OAuthService, sourcesvc *domain.SourceService, group *echo.Group) *Api {
	api := &Api{
		group,
		buildsvc,
		authsvc,
		oauthsvc,
		sourcesvc,
		storage.ImageStorage(),
		storage.UserStorage(),
		storage.BuildStorage(),
		storage.OAuthStorage(),
		storage.AuthStorage(),
	}

	// Authentication endpoints
	group.POST("/auth/login", api.AuthLogin)
	group.POST("/auth/register", api.AuthRegister)

	// Restricted group routes are protected and needs authentication
	restricted := group.Group("", auth.Middleware(authsvc))

	// User endpoints
	restricted.GET("/user/profile", auth.Wrap(api.UserProfile))
	restricted.POST("/user/profile", auth.Wrap(api.UserProfileUpdate))

	// Image endpoints
	restricted.GET("/images", auth.Wrap(api.ImageList))
	restricted.POST("/images", auth.Wrap(api.ImageCreate))
	restricted.GET("/images/:key", auth.Wrap(api.ImageGet))
	restricted.POST("/images/:key", auth.Wrap(api.ImageUpdate))
	restricted.DELETE("/images/:key", auth.Wrap(api.ImageDelete))

	// Build endpoints
	restricted.GET("/builds", auth.Wrap(api.BuildList))
	restricted.GET("/builds/:key", auth.Wrap(api.BuildHistory))

	// OAuth endpoints
	restricted.GET("/oauth/:provider/login_url", auth.Wrap(api.OAuthLogin))
	group.GET("/oauth/:provider/cb/:username", api.OAuthCallback)

	// SourceClient endpoints
	group.POST("/source/:provider/:username/webhook", api.SourceWebhook)
	restricted.GET("/source/:provider/orgs", auth.Wrap(api.SourceOrganisations))
	restricted.GET("/source/:provider/repos", auth.Wrap(api.SourceRepositories))

	return api
}
