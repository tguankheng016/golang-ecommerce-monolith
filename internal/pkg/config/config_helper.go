package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	typeMapper "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/reflection/type_mapper"
)

func BindConfig[T any](e environment.Environment) (T, error) {
	return BindConfigKey[T]("", e)
}

func BindConfigKey[T any](configKey string, env environment.Environment) (T, error) {
	var configPath string

	if env == "" {
		env = environment.ConfigAppEnv()
	}

	// Get default config path
	configPath, err := getConfigRootPath()

	if err != nil {
		return *new(T), err
	}

	// Set the config file name and type
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // type of the config file
	viper.AddConfigPath(configPath) // path to look for the config file

	// Read the default config file
	if err := viper.ReadInConfig(); err != nil {
		return *new(T), err
	}

	// Add the environment-specific config file
	viper.SetConfigName(fmt.Sprintf("config.%s", env)) // look for config.{env}.json

	// Load default config if environment-specific config is not found
	if err := viper.MergeInConfig(); err != nil {
		return *new(T), err
	}

	cfg := typeMapper.GenericInstanceByT[T]()

	// Try to bind
	if err := viper.Unmarshal(&cfg); err != nil {
		return *new(T), err
	}

	return cfg, nil
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
