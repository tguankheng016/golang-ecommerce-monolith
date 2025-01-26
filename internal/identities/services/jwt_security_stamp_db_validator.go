package services

import (
	"context"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/jackc/pgx/v5/pgxpool"
	userService "github.com/tguankheng016/commerce-mono/internal/users/services"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
	"go.uber.org/zap"
)

type securityStampDbValidator struct {
	db           postgres.IPgxDbConn
	cacheManager *cache.Cache[string]
}

func NewSecurityStampDbValidator(db *pgxpool.Pool, cacheManager *cache.Cache[string]) jwt.IJwtSecurityStampDbValidator {
	return &securityStampDbValidator{
		db:           db,
		cacheManager: cacheManager,
	}
}

// ValidateTokenWithStampFromDb checks if a user's token with the given securityStamp is
// valid by checking if there is a matching record in the database and it has not expired.
// If the record is found, then it will cache the security stamp in redis for a certain
// amount of time.
func (c *securityStampDbValidator) ValidateTokenWithStampFromDb(ctx context.Context, cacheKey string, userId int64, securityStamp string) bool {
	userManager := userService.NewUserManager(c.db)

	user, err := userManager.GetUserById(ctx, userId)
	if err != nil {
		logging.Logger.Error("unable to retrieve user", zap.Error(err))
		return false
	}

	if user == nil {
		return false
	}

	if err := c.cacheManager.Set(ctx, cacheKey, user.SecurityStamp.String(), store.WithExpiration(jwt.DefaultCacheExpiration)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in setting cached security stamp", zap.Error(err))
	}

	return user.SecurityStamp.String() == securityStamp
}
