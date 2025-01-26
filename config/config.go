package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/environment"
	"github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
)

type Config struct {
	ServerOptions   *http.ServerOptions       `mapstructure:"serverOptions"`
	PostgresOptions *postgres.PostgresOptions `mapstructure:"postgresOptions"`
	AuthOptions     *jwt.AuthOptions          `mapstructure:"authOptions"`
	RedisOptions    *caching.RedisOptions     `mapstructure:"redisOptions"`
}

func InitConfig(env environment.Environment) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	error,
) {
	var configPath string

	if env == "" {
		env = environment.ConfigAppEnv()
	}

	// Get default config path
	configPath, err := getConfigRootPath()

	if err != nil {
		returnError(err)
	}

	// Set the config file name and type
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // type of the config file
	viper.AddConfigPath(configPath) // path to look for the config file

	// Read the default config file
	if err := viper.ReadInConfig(); err != nil {
		returnError(err)
	}

	// Add the environment-specific config file
	viper.SetConfigName(fmt.Sprintf("config.%s", env)) // look for config.{env}.json

	// Load default config if environment-specific config is not found
	if err := viper.MergeInConfig(); err != nil {
		return returnError(err)
	}

	// Initialize the Config struct and unmarshal the data from Viper into it.
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return returnError(fmt.Errorf("error unmarshaling config, %s", err))
	}

	return &config, config.ServerOptions, config.PostgresOptions, config.AuthOptions, config.RedisOptions, nil
}

func returnError(err error) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	error,
) {
	return nil, nil, nil, nil, nil, err
}

func getConfigRootPath() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up to find the go.mod file
	rootPath := cwd
	for {
		if _, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil {
			// Found the go.mod file
			break
		}
		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			return "", err
		}
		rootPath = parent
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(rootPath, "config")

	return configPath, nil
}
