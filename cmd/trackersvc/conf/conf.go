package conf

import (
	"errors"
)

// Configuration contains necessary information to run apisrv
type Configuration struct {
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

	// JobQ contains configuration for the job queue
	JobQ struct {
		// StatusQueue is the queue name where the build jobs are published
		StatusQueue string `valid:"required"`

		// Driver contains configuration for jobq driver like rabbitmq
		Driver    SingleItemMap          `valid:"required"`
		DriverMap map[string]interface{} `mapstructure:"driver"`
	} `valid:"required"`

	// Storage contains configuration for storage driver like mongodb
	Storage    SingleItemMap          `valid:"required"`
	StorageMap map[string]interface{} `mapstructure:"storage"`
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
