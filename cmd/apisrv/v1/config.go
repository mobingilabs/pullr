package v1

import (
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// APIConfig contains several information necessary to run API endpoints
type APIConfig struct {
	GithubClientID    string
	GithubSecret      string
	ServerURL         string
	FrontendURL       string
	RedirectWhitelist []string
}

// AddConfigFlags adds option flags to command line parser to obtain config
// values from commandline or environment
func AddConfigFlags(set *pflag.FlagSet) {
	set.StringSlice("redirect_whitelist", nil, "Whitelist of urls can be redirected after oauth login requests")
	set.String("github_client", "", "Github client id")
	set.String("github_secret", "", "Github secret")
	set.String("server_url", "", "This server's own url")
	set.String("frontend_url", "", "Frontend url")
}

// ParseConfig reads commandline options and environment variables to populate
// APIConfig.
func ParseConfig() *APIConfig {
	redirectWhitelist := mustList("redirect_whitelist")
	for _, u := range redirectWhitelist {
		mustValidURL(u)
	}

	return &APIConfig{
		GithubClientID:    mustStr("github_client"),
		GithubSecret:      mustStr("github_secret"),
		ServerURL:         mustStr("server_url"),
		FrontendURL:       mustStr("frontend_url"),
		RedirectWhitelist: redirectWhitelist,
	}
}

func mustValidURL(uri string) {
	if _, err := url.Parse(uri); err != nil {
		log.Fatalf("%s is not a valid url", uri)
	}
}

func mustList(key string) []string {
	list := viper.GetStringSlice(key)
	if len(list) == 0 {
		log.Fatalf("%s required to have at least one value", key)
	}

	return list
}

func mustStr(key string) string {
	val := viper.GetString(key)
	if val == "" {
		log.Fatalf("%s required to be set", key)
	}
	return val
}
