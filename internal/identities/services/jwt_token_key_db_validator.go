package services

import (
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
	"go.uber.org/zap"
)

type tokenKeyDbValidator struct {
	db           postgres.IPgxDbConn
	cacheManager *cache.Cache[string]
}

func NewTokenKeyDBValidator(db *pgxpool.Pool, cacheManager *cache.Cache[string]) jwt.IJwtTokenKeyDbValidator {
	return &tokenKeyDbValidator{
		db:           db,
		cacheManager: cacheManager,
	}
}

// ValidateTokenWithTokenKeyFromDb checks if a user's token with the given tokenKey is
// valid by checking if there is a matching record in the database and it has not expired.
// If the record is found, then it will cache the token key in redis for a certain
// amount of time.
func (c *tokenKeyDbValidator) ValidateTokenWithTokenKeyFromDb(ctx context.Context, cacheKey string, userId int64, tokenKey string) bool {
	query := `
		SELECT Count(*) FROM user_tokens 
		WHERE user_id = @userId 
		AND expiration_time > @expirationTime
		AND token_key = @tokenKey 
	`

	args := pgx.NamedArgs{
		"userId":         userId,
		"tokenKey":       tokenKey,
		"expirationTime": time.Now(),
	}

	var count int

	if err := c.db.QueryRow(ctx, query, args).Scan(&count); err != nil {
		logging.Logger.Error("unable to get user token", zap.Error(err))
		return false
	}

	if count == 0 {
		return false
	}

	if err := c.cacheManager.Set(ctx, cacheKey, tokenKey, store.WithExpiration(jwt.DefaultCacheExpiration)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in setting cached token key", zap.Error(err))
	}

	return true
}
