package dummy

import (
	"net/http"

	"github.com/mobingilabs/pullr/pkg/domain"
)

type Github struct {
	Url    string
	Token  string
	Secret string
}

func NewGithub(config domain.OAuthProviderConfig) *Github {
	return &Github{}
}

func (g *Github) LoginUrl(secret string, cbUrl string) string {
	return g.Url
}

func (g *Github) HandleCallback(secret string, req *http.Request) (string, error) {
	return g.Token, nil
}

func (g *Github) GetSecret(req *http.Request) string {
	return g.Secret
}
