package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/eko/gocache/lib/v4/cache"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/tguankheng016/commerce-mono/pkg/caching"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

const (
	SecurityStampKey = "user_security_stamp"
)

type IJwtSecurityStampValidator interface {
	ValidateTokenWithStamp(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type IJwtSecurityStampDbValidator interface {
	ValidateTokenWithStampFromDb(ctx context.Context, cacheKey string, userId int64, securityStamp string) bool
}

type jwtSecurityStampValidator struct {
	dbValidator  IJwtSecurityStampDbValidator
	cacheManager *cache.Cache[string]
}

func NewSecurityStampValidator(dbValidator IJwtSecurityStampDbValidator, cacheManager *cache.Cache[string]) IJwtSecurityStampValidator {
	return &jwtSecurityStampValidator{
		dbValidator:  dbValidator,
		cacheManager: cacheManager,
	}
}

func (j *jwtSecurityStampValidator) ValidateTokenWithStamp(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	securityStamp := claims[SecurityStampKey]
	invalidSecurityStampErr := errors.New("invalid security stamp")

	if securityStamp == nil {
		return invalidSecurityStampErr
	}

	securityStampStr, ok := securityStamp.(string)
	if !ok {
		return invalidSecurityStampErr
	}

	isValid := j.validateTokenWithStampFromCache(ctx, userId, securityStampStr)

	if !isValid {
		isValid = j.dbValidator.ValidateTokenWithStampFromDb(ctx, GenerateStampCacheKey(userId), userId, securityStampStr)
	}

	if !isValid {
		return invalidSecurityStampErr
	}

	return nil
}

func (j *jwtSecurityStampValidator) validateTokenWithStampFromCache(ctx context.Context, userId int64, securityStamp string) bool {
	cacheKey := GenerateStampCacheKey(userId)

	cachedStamp, err := j.cacheManager.Get(ctx, cacheKey)
	if err != nil {
		if !caching.CheckIsCacheValueNotFound(err) {
			logging.Logger.Error("error in getting cached stamp", zap.Error(err))
		}
		return false
	}

	return cachedStamp != "" && cachedStamp == securityStamp
}

func GenerateStampCacheKey(userId int64) string {
	return fmt.Sprintf("%s.%d", SecurityStampKey, userId)
}
