package http

import "fmt"

type ServerOptions struct {
	Port string `mapstructure:"port"                validate:"required" env:"TcpPort"`
	Host string `mapstructure:"host"                                    env:"Host"`
	Name string `mapstructure:"name"                                    env:"ShortTypeName"`
}

func (h *ServerOptions) GetBasePath() string {
	basePath := fmt.Sprintf("http://%s:%s", h.Host, h.Port)

	return basePath
}
