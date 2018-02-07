package v1

import (
	"net/url"
	"os"
	"strings"

	"github.com/golang/glog"
)

type ApiConfig struct {
	GithubClientId    string
	GithubSecret      string
	ServerUrl         string
	FrontendUrl       string
	RedirectWhitelist []string
}

func ParseConfig() *ApiConfig {
	redirectWhitelist := mustEnvList("REDIRECT_WHITELIST")
	for _, u := range redirectWhitelist {
		mustValidUrl(u)
	}

	return &ApiConfig{
		GithubClientId:    mustEnv("GITHUB_CLIENT_ID"),
		GithubSecret:      mustEnv("GITHUB_SECRET"),
		ServerUrl:         mustEnv("SERVER_URL"),
		FrontendUrl:       mustEnv("FRONTEND_URL"),
		RedirectWhitelist: redirectWhitelist,
	}
}

func mustValidUrl(uri string) {
	if _, err := url.Parse(uri); err != nil {
		glog.Fatalf("%s is not a valid url", uri)
	}
}

func mustEnvList(key string) []string {
	val := mustEnv(key)
	list := strings.Split(val, ",")
	if len(list) == 0 {
		glog.Fatalf("%s environment variable required at least one value to run server.", key)
	}

	return list
}

func mustEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		glog.Fatalf("%s environment variable required to run the server.", key)
	}

	return val
}
