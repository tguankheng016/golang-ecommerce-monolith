package constants

const (
	AppEnv                  = "APP_ENV"
	Dev                     = "development"
	Test                    = "test"
	Production              = "production"
	TokenValidityKey        = "token_validity_key"
	RefreshTokenValidityKey = "refresh_token_validity_key"
	SecurityStampKey        = "Identity.SecurityStamp"
	DbContextKey            = "DbContext.Tx"
)

type TxKey string
