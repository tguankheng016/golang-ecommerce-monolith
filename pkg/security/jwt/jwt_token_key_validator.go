package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

const (
	DefaultCacheExpiration  = 1 * time.Hour
	TokenValidityKey        = "token_validity_key"
	RefreshTokenValidityKey = "refresh_token_validity_key"
)

type IJwtTokenKeyValidator interface {
	ValidateTokenWithTokenKey(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type IJwtTokenKeyDbValidator interface {
	ValidateTokenWithTokenKeyFromDb(ctx context.Context, cacheKey string, userId int64, tokenKey string) bool
}

type jwtTokenKeyValidator struct {
	dbValidator  IJwtTokenKeyDbValidator
	cacheManager *cache.Cache[string]
}

func NewTokenKeyValidator(dbValidator IJwtTokenKeyDbValidator, cacheManager *cache.Cache[string]) IJwtTokenKeyValidator {
	return &jwtTokenKeyValidator{
		dbValidator:  dbValidator,
		cacheManager: cacheManager,
	}
}

func (j *jwtTokenKeyValidator) ValidateTokenWithTokenKey(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey := claims[TokenValidityKey]
	invalidTokenKeyErr := errors.New("invalid token key")

	if tokenKey == nil {
		return invalidTokenKeyErr
	}

	tokenKeyStr, ok := tokenKey.(string)
	if !ok {
		return invalidTokenKeyErr
	}

	isValid := j.validateTokenWithTokenKeyFromCache(ctx, userId, tokenKeyStr)

	if !isValid {
		isValid = j.dbValidator.ValidateTokenWithTokenKeyFromDb(ctx, GenerateTokenValidityCacheKey(userId, tokenKeyStr), userId, tokenKeyStr)
	}

	if !isValid {
		return invalidTokenKeyErr
	}

	return nil
}

func (j *jwtTokenKeyValidator) validateTokenWithTokenKeyFromCache(ctx context.Context, userId int64, tokenKey string) bool {
	tokenCacheKey := GenerateTokenValidityCacheKey(userId, tokenKey)

	cachedTokenKey, err := j.cacheManager.Get(ctx, tokenCacheKey)
	if err != nil {
		if !caching.CheckIsCacheValueNotFound(err) {
			logging.Logger.Error("error in getting cached token key", zap.Error(err))
		}
		return false
	}

	return cachedTokenKey != ""
}

func GenerateTokenValidityCacheKey(userId int64, tokenKey string) string {
	return fmt.Sprintf("%s.%d.%s", TokenValidityKey, userId, tokenKey)
}
