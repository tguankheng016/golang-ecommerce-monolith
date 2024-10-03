package config

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/config/environment"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/database"
)

type Config struct {
	GormOptions *database.GormOptions
}

func InitConfig(env environment.Environment) (*Config, *database.GormOptions, error) {
	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, nil, err
	}

	return cfg, cfg.GormOptions, nil
}
