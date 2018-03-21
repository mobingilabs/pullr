package domain_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	. "github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/dummy"
)

type testOAuthProvider struct {
	Test        *testing.T
	WrongSecret bool
	secret      string
}

func (p *testOAuthProvider) LoginURL(secret string, cbUrl string) string {
	p.secret = secret
	return fmt.Sprintf("%s?secret=%s", cbUrl, secret)
}

func (p *testOAuthProvider) FinishLogin(secret string, req *http.Request) (string, error) {
	if p.secret != secret {
		p.Test.Errorf("expected secret: %s, got: %s", p.secret, secret)
	}

	return "testtoken", nil
}

func (t *testOAuthProvider) GetSecret(req *http.Request) string {
	if t.WrongSecret {
		return "somewrongsecret"
	}

	return t.secret
}

func TestOAuthService(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	provider := &testOAuthProvider{Test: t}
	oauthsvc := NewOAuthService(storage.OAuthStorage(), map[string]OAuthProvider{"test": provider})

	_, err := oauthsvc.LoginURL("test", "test", "http://test/api/v1/oauth/github/callback", "redir")
	if err != nil {
		t.Fatal(err)
	}

	token, err := oauthsvc.FinishLogin("test", "test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if token.Token != "testtoken" {
		t.Errorf("expected token: testtoken, got: %s", token)
	}

	if !strings.HasPrefix(token.Redirect, "redir") {
		t.Errorf("expected redirect url to start with: redir, but redirec url is: %s", token.Redirect)
	}
}

func TestOAuthService_BadRequest(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	provider := &testOAuthProvider{Test: t, WrongSecret: true}
	oauthsvc := NewOAuthService(storage.OAuthStorage(), map[string]OAuthProvider{"test": provider})

	_, err := oauthsvc.LoginURL("test", "test", "http://test/api/v1/oauth/github/callback", "redir")
	if err != nil {
		t.Fatal(err)
	}

	_, err = oauthsvc.FinishLogin("test", "test", nil)
	if err != ErrNotFound {
		t.Error("mismatching secrets should result not found error")
	}
}
