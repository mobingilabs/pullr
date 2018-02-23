package conf

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/asaskevich/govalidator.v4"
)

var envPrefix = "PULLR"

func Read() (*Configuration, error) {
	var err error

	// This will allow to override nested configuration values like
	// storage.mongodb.conn as PULLR_STORAGE_MONGODB_CONN
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Configure viper to where to find apisrv configuration files
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")

	viper.SetConfigName("apisrv")
	if confFile := os.Getenv(envPrefix + "_CONF_FILE"); confFile != "" {
		viper.SetConfigFile(confFile)
	}

	viper.SetConfigType("yaml")
	if confType := os.Getenv(envPrefix + "_CONF_TYPE"); confType != "" {
		viper.SetConfigType(confType)
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// set values from env variables starting with PULLR_ to make sure
	// missing keys on config file can be caught over env variables
	envs := os.Environ()
	for _, envKey := range envs {
		keyVal := strings.SplitN(envKey, "=", 2)
		ks := strings.SplitAfterN(keyVal[0], envPrefix+"_", 2)
		if len(ks) != 2 {
			continue
		}

		vKey := strings.ToLower(strings.Replace(ks[1], "_", ".", -1))
		viper.Set(vKey, keyVal[1])
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Auth, err = NewSingleItemMap(config.AuthMap)
	if err != nil {
		return nil, errors.WithMessage(err, "auth config")
	}

	config.Storage, err = NewSingleItemMap(config.StorageMap)
	if err != nil {
		return nil, errors.WithMessage(err, "storage config")
	}

	config.JobQ.Driver, err = NewSingleItemMap(config.JobQ.DriverMap)
	if err != nil {
		return nil, errors.WithMessage(err, "jobq driver config")
	}

	if _, err := govalidator.ValidateStruct(config); err != nil {
		return nil, errors.WithMessage(err, "invalid configuration")
	}

	return &config, nil
}
