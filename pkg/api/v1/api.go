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
	buildsvc *domain.BuildService
	authsvc  *domain.AuthService
	oauthsvc *domain.OAuthService

	// Storages
	imageStorage domain.ImageStorage
	userStorage  domain.UserStorage
	buildStorage domain.BuildStorage
	oauthStorage domain.OAuthStorage
}

// NewApi add api v1 routes to the given routing group
func NewApi(storage domain.StorageDriver, buildsvc *domain.BuildService, authsvc *domain.AuthService, oauthsvc *domain.OAuthService, group *echo.Group) *Api {
	api := &Api{
		group,
		buildsvc,
		authsvc,
		oauthsvc,
		storage.ImageStorage(),
		storage.UserStorage(),
		storage.BuildStorage(),
		storage.OAuthStorage(),
	}

	// Authentication endpoints
	group.POST("/auth/login", api.AuthLogin)
	group.POST("/auth/register", api.AuthRegister)

	// Restricted group routes are protected and needs authentication
	restricted := group.Group("", auth.Middleware(authsvc))

	// User endpoints
	restricted.GET("/user/profile", auth.Wrap(api.UserProfile))
	restricted.GET("/user/profile", auth.Wrap(api.UserProfileUpdate))

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
	restricted.GET("/oauth/:provider/:username/cb", api.OAuthCallback)

	// SourceClient endpoints
	group.POST("/:provider/webhook", api.SourceWebhook)
	restricted.GET("/:provider/organisations", auth.Wrap(api.SourceOrganisations))
	restricted.GET("/:provider/repositories", auth.Wrap(api.SourceRepositories))

	return api
}
