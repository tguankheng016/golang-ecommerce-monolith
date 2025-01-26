package environment

import "os"

type Environment string

var (
	Development = Environment("development")
	Test        = Environment("test")
	Production  = Environment("production")
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

func (env Environment) IsTest() bool {
	return env == Test
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}
