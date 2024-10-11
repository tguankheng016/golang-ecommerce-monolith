package jwt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

type IJwtTokenValidator interface {
	ValidateToken(ctx context.Context, token string, tokenType TokenType) (int64, jwtGo.MapClaims, error)
}

type jwtTokenValidator struct {
	secretKey string
	issuer    string
	audience  string
	db        *gorm.DB
	client    *redis.Client
	logger    logger.ILogger
}

const (
	DefaultCacheExpiration = 1 * time.Hour
)

func NewJwtTokenValidator(db *gorm.DB, client *redis.Client, logger logger.ILogger, authOptions *AuthOptions) IJwtTokenValidator {
	return &jwtTokenValidator{
		secretKey: authOptions.SecretKey,
		issuer:    authOptions.Issuer,
		audience:  authOptions.Audience,
		db:        db,
		client:    client,
		logger:    logger,
	}
}

func (j *jwtTokenValidator) ValidateToken(ctx context.Context, tokenString string, tokenType TokenType) (int64, jwtGo.MapClaims, error) {
	token, err := jwtGo.ParseWithClaims(tokenString, jwtGo.MapClaims{}, func(token *jwtGo.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return 0, nil, err
	}

	if claims, ok := token.Claims.(jwtGo.MapClaims); ok && token.Valid {
		// validate token type
		tokenTypeInt, _ := strconv.Atoi(claims["token_type"].(string))

		if tokenTypeInt != int(tokenType) {
			return 0, nil, errors.New("Invalid token type")
		}

		// token is valid and has not expired
		iss := token.Header["iss"]
		if iss != j.issuer {
			// handle invalid issuer
			return 0, nil, errors.New("Invalid token issuer")
		}

		aud := token.Header["aud"]
		if aud != j.audience {
			// handle invalid audience
			return 0, nil, errors.New("Invalid token audience")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			// handle error
			return 0, nil, errors.New("Invalid sub")
		}

		userId, err := strconv.ParseInt(sub, 10, 64)

		if err != nil {
			return 0, nil, err
		}

		if err := j.validateTokenWithSecurityStamp(ctx, userId, claims); err != nil {
			return 0, nil, err
		}

		if err := j.validateTokenWithTokenKey(ctx, userId, claims); err != nil {
			return 0, nil, err
		}

		return userId, claims, nil
	}

	return 0, nil, errors.New("Invalid token")
}

func (j *jwtTokenValidator) validateTokenWithSecurityStamp(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	securityStamp := claims[constants.SecurityStampKey]
	invalidSecurityStampErr := errors.New("Invalid stamp")

	if securityStamp == nil {
		return invalidSecurityStampErr
	}

	isValid := j.validateTokenWithSecurityStampFromCache(ctx, userId, securityStamp.(string))

	if !isValid {
		isValid = j.validateTokenWithSecurityStampFromDb(ctx, userId, securityStamp.(string))
	}

	if !isValid {
		return invalidSecurityStampErr
	}

	return nil
}

func (j *jwtTokenValidator) validateTokenWithSecurityStampFromCache(ctx context.Context, userId int64, securityStamp string) bool {
	cacheKey := generateStampCacheKey(userId)

	cachedStamp, err := j.client.Get(ctx, cacheKey).Result()
	if err != nil {
		return false
	}

	return cachedStamp != "" && cachedStamp == securityStamp
}

func (j *jwtTokenValidator) validateTokenWithSecurityStampFromDb(ctx context.Context, userId int64, securityStamp string) bool {
	var user models.User
	if err := j.db.First(&user, userId).Error; err != nil {
		return false
	}

	if err := j.client.Set(ctx, generateStampCacheKey(userId), user.SecurityStamp.String(), DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		j.logger.Error(err)
	}

	if user.SecurityStamp.String() != securityStamp {
		return false
	}

	return true
}

func (j *jwtTokenValidator) validateTokenWithTokenKey(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey := claims[constants.TokenValidityKey]
	invalidTokenKeyErr := errors.New("Invalid token key")

	if tokenKey == nil {
		return invalidTokenKeyErr
	}

	isValid := j.validateTokenWithTokenKeyFromCache(ctx, userId, tokenKey.(string))

	if !isValid {
		isValid = j.validateTokenWithTokenKeyFromDb(ctx, userId, tokenKey.(string))
	}

	if !isValid {
		return invalidTokenKeyErr
	}

	return nil
}

func (j *jwtTokenValidator) validateTokenWithTokenKeyFromCache(ctx context.Context, userId int64, tokenKey string) bool {
	tokenCacheKey := generateTokenValidityCacheKey(userId, tokenKey)

	cachedTokenKey, err := j.client.Get(ctx, tokenCacheKey).Result()
	if err != nil {
		return false
	}

	return cachedTokenKey != ""
}

func (j *jwtTokenValidator) validateTokenWithTokenKeyFromDb(ctx context.Context, userId int64, tokenKey string) bool {
	tokenCacheKey := generateTokenValidityCacheKey(userId, tokenKey)

	var count int64
	if err := j.db.Model(&models.UserToken{}).Where("user_id = ? AND token_key = ? AND expiration_time > ?", userId, tokenKey, time.Now()).Count(&count).Error; err != nil || count == 0 {
		return false
	}

	if err := j.client.Set(ctx, tokenCacheKey, tokenKey, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		j.logger.Error(err)
	}

	return true
}

func generateStampCacheKey(userId int64) string {
	return fmt.Sprintf("%s.%d", constants.SecurityStampKey, userId)
}

func generateTokenValidityCacheKey(userId int64, tokenKey string) string {
	return fmt.Sprintf("%s.%d.%s", constants.TokenValidityKey, userId, tokenKey)
}
