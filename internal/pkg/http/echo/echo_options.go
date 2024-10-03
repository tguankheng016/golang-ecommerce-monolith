package echoserver

import (
	"fmt"
	"net/url"
)

type EchoOptions struct {
	Port                string   `mapstructure:"port"                validate:"required" env:"TcpPort"`
	Development         bool     `mapstructure:"development"                             env:"Development"`
	BasePath            string   `mapstructure:"basePath"            validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"                     env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"                                 env:"Timeout"`
	Host                string   `mapstructure:"host"                                    env:"Host"`
	Name                string   `mapstructure:"name"                                    env:"ShortTypeName"`
}

func (c *EchoOptions) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c *EchoOptions) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}
