package config

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
)

type Config struct {
	GormOptions *database.GormOptions   `mapstructure:"gormOptions"`
	EchoOptions *echoServer.EchoOptions `mapstructure:"echoOptions"`
	Logger      *logger.LoggerConfig    `mapstructure:"logOptions"`
}

func InitConfig(env environment.Environment) (*Config, *database.GormOptions,
	*echoServer.EchoOptions, *logger.LoggerConfig, error) {

	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return cfg, cfg.GormOptions, cfg.EchoOptions, cfg.Logger, nil
}
