package redis

type RedisOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database int    `mapstructure:"database"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"poolSize"`
}
