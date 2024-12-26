package jwt

type AuthOptions struct {
	SecretKey string `mapstructure:"secretKey"`
	Issuer    string `mapstructure:"issuer"`
	Audience  string `mapstructure:"audience"`
}
