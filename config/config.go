package config

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/redis"
)

type Config struct {
	GormOptions  *database.GormOptions   `mapstructure:"gormOptions"`
	EchoOptions  *echoServer.EchoOptions `mapstructure:"echoOptions"`
	Logger       *logger.LoggerConfig    `mapstructure:"logOptions"`
	AuthOptions  *jwt.AuthOptions        `mapstructure:"authOptions"`
	RedisOptions *redis.RedisOptions     `mapstructure:"redisOptions"`
}

func InitConfig(env environment.Environment) (*Config, *database.GormOptions,
	*echoServer.EchoOptions, *logger.LoggerConfig, *jwt.AuthOptions, *redis.RedisOptions, error) {

	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return cfg, cfg.GormOptions, cfg.EchoOptions, cfg.Logger, cfg.AuthOptions, cfg.RedisOptions, nil
}
