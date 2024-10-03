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
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fmt.Println(fmt.Sprintf("Current working directory is: %s", wd))

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")

	return configPath, nil
}
