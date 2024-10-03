package environment

import (
	"os"

	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
)

type Environment string

var (
	Development = Environment(constants.Dev)
	Test        = Environment(constants.Test)
	Production  = Environment(constants.Production)
)

func ConfigAppEnv() Environment {
	environment := Development

	// Find environment variable for the environment
	env := os.Getenv("APP_ENV")

	if env != "" {
		environment = Environment(env)
	}

	return environment
}

func (env Environment) IsDevelopment() bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}
