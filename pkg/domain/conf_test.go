package domain

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var confFixture string

var expectedConf = Config{
	Log:      LogConfig{"info", "text"},
	OAuth:    map[string]OAuthProviderConfig{"github": {"id", "secret"}},
	ApiSrv:   ApiSrvConfig{false, []string{"*"}, 8080},
	BuildCtl: BuildCtlConfig{"pullr-image-build", 5, "./src", time.Minute * 5},
	Storage: DriverConfig{
		Driver: "mongodb",
		Options: map[string]interface{}{
			"conn": "mongodb://pullr:pullrpass@pullr-mongodb/pullr",
		},
	},
	JobQ: DriverConfig{
		Driver: "rabbitmq",
		Options: map[string]interface{}{
			"conn": "amqp://pullr-rabbitmq:5672",
		},
	},
	Auth: AuthConfig{
		Key: "/certs/auth.key",
		Crt: "/certs/auth.crt",
	},
	Registry: RegistryConfig{
		URL:      "https://docker-registry:5050",
		Username: "user",
		Password: "pass",
	},
}

func TestMain(m *testing.M) {
	pullrYaml, err := ioutil.ReadFile("../../conf/pullr.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read conf/pullr.yml, it is used as test fixture: %v", err)
		os.Exit(1)
	}

	confFixture = string(pullrYaml)
	os.Exit(m.Run())
}

func TestParseConfig(t *testing.T) {
	reader := strings.NewReader(confFixture)
	conf, err := ParseConfig(reader)
	if err != nil {
		t.Error(err)
	}

	assertConf(t, conf)
}

func TestConfig_SetByEnv(t *testing.T) {
	reader := strings.NewReader(confFixture)
	conf, err := ParseConfig(reader)
	if err != nil {
		t.Error(err)
	}

	env := []string{
		"PULLR_NOKEY=noval",
		"PULLR_APISRV_PORT=9090",
		"PULLR_JOBQ_OPTIONS_NEWKEY=newval",
		"JOBQ_OPTIONS_IGNORE=ignored",
	}

	conf.SetByEnv("PULLR", env)

	assertEq(t, "apisrv.http", conf.ApiSrv.Port, 9090)
	assertEq(t, "jobq.options.newkey", conf.JobQ.Options["newkey"], "newval")
	_, ok := conf.JobQ.Options["ignore"]
	assert(t, !ok, "env variables without matching prefix should be ignored")
}

func assertConf(t *testing.T, conf *Config) {
	expected := expectedConf
	assert(t, conf != nil, "config parsed successfully but conf pointer is nil")

	assertEq(t, "log.level", conf.Log.Level, "info")
	assertEq(t, "log.formatter", conf.Log.Formatter, "text")

	assertEq(t, "apisrv.enableCors", conf.ApiSrv.EnableCORS, expected.ApiSrv.EnableCORS)
	assertEq(t, "apisrv.port", conf.ApiSrv.Port, expected.ApiSrv.Port)
	assertEq(t, "apisrv.alloworigins.length", len(conf.ApiSrv.AllowOrigins), len(expected.ApiSrv.AllowOrigins))
	for i := range expected.ApiSrv.AllowOrigins {
		assertEq(t, fmt.Sprintf("apisrv.alloworigins.%d", i), conf.ApiSrv.AllowOrigins[i], expected.ApiSrv.AllowOrigins[i])
	}

	assertEq(t, "buildctl", conf.BuildCtl, expected.BuildCtl)

	assertEq(t, "oauth.github", conf.OAuth["github"], expected.OAuth["github"])

	assertEq(t, "storage", conf.Storage.Driver, expected.Storage.Driver)
	assertEqMap(t, "storage.mongodb", conf.Storage.Options, expected.Storage.Options)

	assertEq(t, "jobq", conf.JobQ.Driver, expected.JobQ.Driver)
	assertEqMap(t, "jobq.rabbitmq", conf.JobQ.Options, expected.JobQ.Options)

	assertEq(t, "auth", conf.Auth, expected.Auth)
	assertEq(t, "registry", conf.Registry, expected.Registry)
}

func assert(t *testing.T, assertion bool, err string) {
	if !assertion {
		t.Error(err)
	}
}

func assertEq(t *testing.T, path string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%s expected: %v actual: %v", path, expected, actual)
	}
}

func assertEqMap(t *testing.T, path string, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s expectd: %v actual: %v", path, expected, actual)
	}
}
