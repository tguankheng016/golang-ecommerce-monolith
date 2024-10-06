package jwt

import (
	"strconv"
	"time"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"gorm.io/gorm"
)

type IJwtTokenValidator interface {
	ValidateToken(token string, tokenType TokenType) error
}

type jwtTokenValidator struct {
	secretKey string
	issuer    string
	audience  string
	db        *gorm.DB
}

func NewJwtTokenValidator(db *gorm.DB, authOptions *AuthOptions) IJwtTokenValidator {
	return &jwtTokenValidator{
		secretKey: authOptions.SecretKey,
		issuer:    authOptions.Issuer,
		audience:  authOptions.Audience,
		db:        db,
	}
}

func (j *jwtTokenValidator) ValidateToken(tokenString string, tokenType TokenType) error {
	token, err := jwtGo.ParseWithClaims(tokenString, jwtGo.MapClaims{}, func(token *jwtGo.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		// validate token type
		tokenTypeInt, _ := strconv.Atoi(claims["token_type"].(string))

		if tokenTypeInt != int(tokenType) {
			return errors.New("Invalid token type")
		}

		// token is valid and has not expired
		iss := token.Header["iss"]
		if iss != j.issuer {
			// handle invalid issuer
			return errors.New("Invalid token issuer")
		}

		aud := token.Header["aud"]
		if aud != j.audience {
			// handle invalid audience
			return errors.New("Invalid token audience")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			// handle error
			return errors.New("Invalid sub")
		}

		userId, err := strconv.ParseInt(sub, 10, 64)

		if err != nil {
			return err
		}

		if err := j.validateTokenWithSecurityStamp(userId, claims); err != nil {
			return err
		}

		if err := j.validateTokenWithTokenKey(userId, claims); err != nil {
			return err
		}

	}

	return nil
}

func (j *jwtTokenValidator) validateTokenWithSecurityStamp(userId int64, claims jwtGo.MapClaims) error {
	securityStamp := claims[constants.SecurityStampKey]
	invalidSecurityStampErr := errors.New("Invalid stamp")

	if securityStamp == nil {
		return invalidSecurityStampErr
	}

	var user models.User
	if err := j.db.First(&user, userId).Error; err != nil {
		return err
	}

	if user.SecurityStamp.String() != securityStamp {
		return invalidSecurityStampErr
	}

	return nil
}

func (j *jwtTokenValidator) validateTokenWithTokenKey(userId int64, claims jwtGo.MapClaims) error {
	tokenKey := claims[constants.TokenValidityKey]
	invalidTokenKeyErr := errors.New("Invalid token key")

	if tokenKey == nil {
		return invalidTokenKeyErr
	}

	var count int64
	if err := j.db.Model(&models.UserToken{}).Where("user_id = ? AND token_key = ?", userId, tokenKey).Where("expiration_time > ?", time.Now()).Count(&count).Error; err != nil || count == 0 {
		return invalidTokenKeyErr
	}

	return nil
}
