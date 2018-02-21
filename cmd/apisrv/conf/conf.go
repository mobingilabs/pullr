package conf

import (
	"errors"
)

// Configuration contains necessary information to run apisrv
type Configuration struct {
	// Auth contains configuration for the pullr's authorization service driver
	Auth    SingleItemMap          `yaml:"-" valid:"required"`
	AuthMap map[string]interface{} `mapstructure:"auth"`

	// Storage contains configuration for the pullr's storage service driver
	Storage    SingleItemMap          `valid:"required"`
	StorageMap map[string]interface{} `mapstructure:"storage"`

	// WebhookURL is the url of the VCS webhook callback handler endpoint
	WebhookURL string `valid:"required"`

	// Log contains configuration parameters for logging
	Log struct {
		// ForceColors is tells logger to use always colored output
		ForceColors bool

		// Level is the level at which operations are logged.
		// This can be error, warn, info, or debug.
		Level string `valid:"in(error|warn|info|debug),required"`

		// Formatter overrides the default formatter with another. Options
		// include "text", "json" and "logstash".
		Formatter string `valid:"in(text|json),required"`
	}

	// HTTP contains configuration parameters for the api server
	HTTP struct {
		// EnableCORS enables cross origin requests
		EnableCORS bool

		// AllowOrigins is the list of trusted addresses to use with CORS
		AllowOrigins []string

		// Port number to listen for incoming requests
		Port int `valid:"required"`
	} `valid:"required"`

	// OAuth contains configuration parameters for getting user's source code
	// information
	OAuth struct {
		// RedirectWhitelist is a whitelist of urls which they can be redirected
		// to after successful oauth logins
		RedirectWhitelist []string `valid:"required"`

		// CallbackURL will be called after users logged in to oauth provider
		CallbackURL string `valid:"required"`

		// Clients is a dictionary of oauth client configurations
		Clients map[string]struct {
			// ID is the client id given by the oauth provider
			ID string `valid:"required"`

			// Secret is the client secret given by the oauth provider
			Secret string `valid:"required"`
		} `valid:"required"`
	}
}

type SingleItemMap struct {
	Name       string
	Parameters map[string]interface{}
}

func NewSingleItemMap(in map[string]interface{}) (SingleItemMap, error) {
	var sim SingleItemMap

	if len(in) == 0 {
		return sim, errors.New("exactly one configuration required")
	}

	for childKey, params := range in {
		sim.Name = childKey
		sim.Parameters = params.(map[string]interface{})
		break
	}

	// Make sure env variables evaluated
	//readDynamicConfTree([]string{key, sim.Name}, sim.Parameters)
	return sim, nil
}
