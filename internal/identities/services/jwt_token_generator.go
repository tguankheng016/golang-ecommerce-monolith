package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/gofrs/uuid"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"github.com/tguankheng016/commerce-mono/pkg/postgres"
	"github.com/tguankheng016/commerce-mono/pkg/security/jwt"
	"go.uber.org/zap"
)

const (
	AccessTokenExpirationTime  = 24 * time.Hour
	RefreshTokenExpirationTime = 30 * 24 * time.Hour
)

type IJwtTokenGenerator interface {
	GenerateAccessToken(ctx context.Context, user models.User, refreshTokenKey string) (string, int, error)
	GenerateRefreshToken(ctx context.Context, user models.User) (string, string, int, error)
	RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type jwtTokenGenerator struct {
	secretKey    string
	issuer       string
	audience     string
	db           postgres.IPgxDbConn
	cacheManager *cache.Cache[string]
}

func NewJwtTokenGenerator(db *pgxpool.Pool, cacheManager *cache.Cache[string], authOptions *jwt.AuthOptions) IJwtTokenGenerator {
	return &jwtTokenGenerator{
		secretKey:    authOptions.SecretKey,
		issuer:       authOptions.Issuer,
		audience:     authOptions.Audience,
		db:           db,
		cacheManager: cacheManager,
	}
}

// GenerateAccessToken generates an access token for the given user.
//
// The function takes in a user model and the refresh token key, and returns a JWT token, the expiration time in seconds, and an error.
// The expiration time is set to the AccessTokenExpirationTime duration.
// The function will also insert a new row into the user_tokens table with the user's ID, the token key, and the expiration time.
// Finally, the function will cache the token validity key in Redis with the key "user_token:<user_id>:<token_key>" and the expiration time set to the DefaultCacheExpiration.
// If there is an error caching the token validity key, the function will log the error but not return it.
// If there is an error, the function will return an empty string, 0, and the error.
func (j *jwtTokenGenerator) GenerateAccessToken(ctx context.Context, user models.User, refreshTokenKey string) (string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, jwt.AccessToken, refreshTokenKey)

	if err != nil {
		return "", 0, err
	}

	accessToken, err := j.createToken(claims)

	return accessToken, int(AccessTokenExpirationTime.Seconds()), err
}

// GenerateRefreshToken generates a refresh token for the given user.
// The function takes in a user model and returns a JWT token, the token key, and the expiration time in seconds.
// If there is an error, the function will return an empty string, an empty string, 0, and the error.
func (j *jwtTokenGenerator) GenerateRefreshToken(ctx context.Context, user models.User) (string, string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, jwt.RefreshToken, "")

	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := j.createToken(claims)

	refreshTokenKey := claims[jwt.TokenValidityKey]
	refreshTokenStr := fmt.Sprintf("%s", refreshTokenKey)

	return refreshToken, refreshTokenStr, int(RefreshTokenExpirationTime.Seconds()), err
}

// RemoveUserTokens removes a user's access token and refresh token from the database and Redis cache.
//
// It takes in the user's ID and the claims map from the JWT token. It checks if the claims map contains
// the token validity key and removes the token from the database and Redis cache. If the claims map also
// contains the refresh token validity key, it removes the refresh token from the database and Redis cache
// as well.
//
// If there is an error removing the tokens, it returns the error.
func (j *jwtTokenGenerator) RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey, ok := claims[jwt.TokenValidityKey]
	if !ok {
		return errors.New("invalid token key")
	}

	if err := j.removeToken(ctx, userId, tokenKey.(string)); err != nil {
		return err
	}

	refreshTokenKey, ok := claims[jwt.RefreshTokenValidityKey]
	if ok {
		if err := j.removeToken(ctx, userId, refreshTokenKey.(string)); err != nil {
			return err
		}
	}

	return nil
}

// createToken creates a new JWT token with the given claims.
//
// The function takes in a MapClaims object and returns a signed string
// representation of the token. The signing method used is HS256.
//
// The function also sets the "iss" and "aud" headers of the token to the
// issuer and audience of the JWT, respectively.
//
// If there is an error signing the token, the function will return the error.
func (j *jwtTokenGenerator) createToken(claims jwtGo.MapClaims) (string, error) {
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token.Header["iss"] = j.issuer
	token.Header["aud"] = j.audience

	return token.SignedString([]byte(j.secretKey))
}

// createJwtClaims creates a new JWT claims map from the given user and token type.
//
// The claims map will contain the user's ID, the JWT ID, the issued at time, the not before time,
// the expiration time, the token validity key, the security stamp, and the token type.
//
// If the token type is RefreshToken, the expiration time will be set to the RefreshTokenExpirationTime.
// Otherwise, the expiration time will be set to the AccessTokenExpirationTime.
//
// The function will also insert a new row into the user_tokens table with the user's ID, the token key,
// and the expiration time.
//
// Finally, the function will cache the token validity key in Redis with the key
// "user_token:<user_id>:<token_key>" and the expiration time set to the DefaultCacheExpiration.
// If there is an error caching the token validity key, the function will log the error but not return
// it.
//
// Returns the claims map and an error if there is one.
func (j *jwtTokenGenerator) createJwtClaims(ctx context.Context, user models.User, tokenType jwt.TokenType, refreshTokenKey string) (jwtGo.MapClaims, error) {
	tokenValidityKey, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	claims := jwtGo.MapClaims{}

	claims["jti"], err = uuid.NewV6()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var expiration time.Duration
	if tokenType == jwt.RefreshToken {
		expiration = RefreshTokenExpirationTime
	} else {
		expiration = AccessTokenExpirationTime
	}

	claims["sub"] = strconv.FormatInt(user.Id, 10)
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(expiration).Unix()
	claims[jwt.TokenValidityKey] = tokenValidityKey
	claims[jwt.SecurityStampKey] = user.SecurityStamp
	claims["token_type"] = strconv.Itoa(int(tokenType))

	if refreshTokenKey != "" {
		claims[jwt.RefreshTokenValidityKey] = refreshTokenKey
	}

	// Add User Token
	query := `
		INSERT INTO user_tokens (user_id, token_key, expiration_time)
		VALUES (@userId, @tokenKey, @expirationTime)
	`

	args := pgx.NamedArgs{
		"userId":         user.Id,
		"tokenKey":       tokenValidityKey.String(),
		"expirationTime": now.Add(expiration),
	}

	if _, err := j.db.Exec(ctx, query, args); err != nil {
		return nil, err
	}

	if err := j.cacheManager.Set(ctx, jwt.GenerateTokenValidityCacheKey(user.Id, tokenValidityKey.String()), tokenValidityKey.String(), store.WithExpiration(jwt.DefaultCacheExpiration)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in setting cached token key", zap.Error(err))
	}

	return claims, nil
}

// removeToken deletes a user's token with the given tokenKey from the database and the cache.
// It is used to invalidate a user's token, which is necessary for logging out a user.
func (j *jwtTokenGenerator) removeToken(ctx context.Context, userId int64, tokenKey string) error {
	query := "DELETE FROM user_tokens WHERE user_id = @userId and token_key = @tokenKey"

	args := pgx.NamedArgs{
		"userId":   userId,
		"tokenKey": tokenKey,
	}

	if _, err := j.db.Exec(ctx, query, args); err != nil {
		return err
	}

	if err := j.cacheManager.Delete(ctx, jwt.GenerateTokenValidityCacheKey(userId, tokenKey)); err != nil {
		// Dont return just log
		logging.Logger.Error("error in deleting cached token key", zap.Error(err))
	}

	return nil
}
