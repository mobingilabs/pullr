package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// UserProfile is an handler serves only authenticated users. Responds with
// authenticated user's profile data
func (a *Api) UserProfile(secrets domain.AuthSecrets, c echo.Context) error {
	type response struct {
		User   domain.User `json:"user"`
		Tokens []string    `json:"tokens"`
	}

	usr, err := a.userStorage.Get(secrets.Username)
	if err != nil {
		return err
	}

	tokens, err := a.oauthStorage.GetTokens(secrets.Username)
	var tokenProviders []string
	for provider := range tokens {
		tokenProviders = append(tokenProviders, provider)
	}

	return c.JSON(http.StatusOK, response{usr, tokenProviders})
}

// UserProfileUpdate is an handle
func (a *Api) UserProfileUpdate(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
