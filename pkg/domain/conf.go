package domain

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
)

// Config is the configuration object for pullr
type Config struct {
	Storage DriverConfig
	JobQ    DriverConfig

	Auth AuthConfig
	Log  LogConfig `valid:"-"`

	BuildCtl BuildCtlConfig
	ApiSrv   ApiSrvConfig

	OAuth map[string]OAuthProviderConfig

	Registry RegistryConfig
}

// DriverConfig is a pair of an implementation name and a configuration for that
// implementation for the services
type DriverConfig struct {
	Driver  string
	Options map[string]interface{}
}

// AuthConfig contains configuration options for authentication service
type AuthConfig struct {
	Key string
	Crt string
}

// LogConfig contains configuration for logger
type LogConfig struct {
	Level     string `valid:"in(error|warn|info|debug)"`
	Formatter string `valid:"in(text|json)"`
}

// BuildCtlConfig contains configuration for buildctl service
type BuildCtlConfig struct {
	Queue    string
	MaxErr   int
	CloneDir string
	Timeout  time.Duration
}

// ApiSrvConfig contains configuration for apisrv service
type ApiSrvConfig struct {
	EnableCORS   bool     `valid:"-"`
	AllowOrigins []string `valid:"-"`
	Port         int
}

// OAuthProviderConfig is configuration for authenticating with oauth providers
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
}

// RegistryConfig contains configuration for a docker registry to push images
type RegistryConfig struct {
	URL      string
	Username string
	Password string
}

// ParseConfig parses given yaml/json input into Config
func ParseConfig(reader io.Reader) (*Config, error) {
	var conf Config
	err := yaml.NewDecoder(reader).Decode(&conf)
	return &conf, err
}

// SetByEnv overrides the configuration by the environment variables
func (c *Config) SetByEnv(prefix string, keyValues []string) {
	for _, val := range keyValues {
		if !strings.HasPrefix(val, prefix) {
			continue
		}

		env := strings.SplitN(val, "=", 2)
		fieldPath := strings.Split(env[0], "_")

		if len(env) < 2 || len(fieldPath) < 2 {
			continue
		}

		// Skip PREFIX_
		fieldPath = fieldPath[1:]
		for i := range fieldPath {
			fieldPath[i] = strings.ToLower(fieldPath[i])
		}

		err := c.setField(fieldPath, reflect.ValueOf(c), env[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: couldn't set configuration from env: %v\n", err)
		}
	}
}

func (c *Config) setField(path []string, node reflect.Value, value string) error {
	switch node.Kind() {
	case reflect.Ptr:
		node = node.Elem()
		if !node.IsValid() {
			node.Set(reflect.New(node.Type()))
		}
		return c.setField(path, node, value)

	case reflect.Interface:
		node = node.Elem()
		return c.setField(path, node, value)

	case reflect.Struct:
		for i := 0; i < node.NumField(); i++ {
			nameLower := strings.ToLower(node.Type().Field(i).Name)
			if nameLower == path[0] {
				node = node.Field(i)
				return c.setField(path[1:], node, value)
			}
		}

		return fmt.Errorf("field not found: %s", path[0])

	case reflect.Map:
		for _, key := range node.MapKeys() {
			if strings.ToLower(key.String()) == path[0] {
				elem := node.MapIndex(key)
				copy := reflect.New(elem.Type()).Interface()
				if err := copier.Copy(copy, elem.Interface()); err != nil {
					return err
				}

				err := c.setField(path[1:], reflect.ValueOf(copy), value)
				if err != nil {
					return err
				}

				node.SetMapIndex(key, reflect.ValueOf(copy).Elem())
				return nil
			}
		}

		key := reflect.ValueOf(path[0])
		switch node.Type().Elem().Kind() {
		case reflect.String, reflect.Interface:
			node.SetMapIndex(key, reflect.ValueOf(value))
			return nil
		case reflect.Bool:
			lower := strings.ToLower(value)
			truth := lower == "1" || value == "true"
			node.SetMapIndex(key, reflect.ValueOf(truth))
		case reflect.Int:
			intVal, err := strconv.ParseInt(value, 10, strconv.IntSize)
			if err != nil {
				return err
			}
			node.SetMapIndex(key, reflect.ValueOf(intVal))
			return nil
		}

	case reflect.String:
		node.SetString(value)
		return nil

	case reflect.Int:
		intVal, err := strconv.ParseInt(value, 10, strconv.IntSize)
		if err != nil {
			return err
		}

		node.SetInt(intVal)
		return nil
	}

	return fmt.Errorf("unsupported type: %v", node.Kind())
}
