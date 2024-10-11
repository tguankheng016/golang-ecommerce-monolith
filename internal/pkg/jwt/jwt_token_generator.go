package jwt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
	"gorm.io/gorm"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

const (
	AccessTokenExpirationTime  = 24 * time.Hour
	RefreshTokenExpirationTime = 30 * 24 * time.Hour
)

type IJwtTokenGenerator interface {
	GenerateAccessToken(ctx context.Context, user *models.User, refreshTokenKey string) (string, int, error)
	GenerateRefreshToken(ctx context.Context, user *models.User) (string, string, int, error)
	RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error
}

type jwtTokenGenerator struct {
	secretKey string
	issuer    string
	audience  string
	db        *gorm.DB
	client    *redis.Client
	logger    logger.ILogger
}

func NewJwtTokenGenerator(db *gorm.DB, client *redis.Client, logger logger.ILogger, authOptions *AuthOptions) IJwtTokenGenerator {
	return &jwtTokenGenerator{
		secretKey: authOptions.SecretKey,
		issuer:    authOptions.Issuer,
		audience:  authOptions.Audience,
		db:        db,
		client:    client,
		logger:    logger,
	}
}

func (j *jwtTokenGenerator) GenerateAccessToken(ctx context.Context, user *models.User, refreshTokenKey string) (string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, AccessToken, refreshTokenKey)

	if err != nil {
		return "", 0, err
	}

	accessToken, err := j.createToken(claims)

	return accessToken, int(AccessTokenExpirationTime.Seconds()), err
}

func (j *jwtTokenGenerator) GenerateRefreshToken(ctx context.Context, user *models.User) (string, string, int, error) {
	claims, err := j.createJwtClaims(ctx, user, RefreshToken, "")

	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := j.createToken(claims)

	refreshTokenKey := claims[constants.TokenValidityKey]
	refreshTokenStr := fmt.Sprintf("%s", refreshTokenKey)

	return refreshToken, refreshTokenStr, int(RefreshTokenExpirationTime.Seconds()), err
}

func (j *jwtTokenGenerator) RemoveUserTokens(ctx context.Context, userId int64, claims jwtGo.MapClaims) error {
	tokenKey, ok := claims[constants.TokenValidityKey]
	if !ok {
		return errors.New("Invalid token key")
	}

	if err := j.removeToken(ctx, userId, tokenKey.(string)); err != nil {
		return err
	}

	refreshTokenKey, ok := claims[constants.RefreshTokenValidityKey]
	if ok {
		if err := j.removeToken(ctx, userId, refreshTokenKey.(string)); err != nil {
			return err
		}
	}

	return nil
}

func (j *jwtTokenGenerator) createToken(claims jwtGo.MapClaims) (string, error) {
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token.Header["iss"] = j.issuer
	token.Header["aud"] = j.audience

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtTokenGenerator) createJwtClaims(ctx context.Context, user *models.User, tokenType TokenType, refreshTokenKey string) (jwtGo.MapClaims, error) {
	tokenValidityKey, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	claims := jwtGo.MapClaims{}

	claims["jti"], err = uuid.NewV4()

	if err != nil {
		return nil, err
	}

	now := time.Now()
	var expiration time.Duration
	if tokenType == RefreshToken {
		expiration = RefreshTokenExpirationTime
	} else {
		expiration = AccessTokenExpirationTime
	}

	claims["sub"] = strconv.FormatInt(user.Id, 10)
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(expiration).Unix()
	claims[constants.TokenValidityKey] = tokenValidityKey
	claims[constants.SecurityStampKey] = user.SecurityStamp
	claims["token_type"] = strconv.Itoa(int(tokenType))

	if refreshTokenKey != "" {
		claims[constants.RefreshTokenValidityKey] = refreshTokenKey
	}

	// Add User Token
	userToken := models.UserToken{
		UserId:         user.Id,
		TokenKey:       tokenValidityKey.String(),
		ExpirationTime: now.Add(expiration),
	}

	if err := j.db.Create(&userToken).Error; err != nil {
		return nil, errors.Wrap(err, "error when inserting user token into the database.")
	}

	if err := j.client.Set(ctx, generateTokenValidityCacheKey(user.Id, tokenValidityKey.String()), tokenValidityKey.String(), DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		j.logger.Error(err)
	}

	return claims, nil
}

func (j *jwtTokenGenerator) removeToken(ctx context.Context, userId int64, tokenKey string) error {
	if err := j.db.Where("token_key = ? AND user_id = ?", tokenKey, userId).Delete(&models.UserToken{}).Error; err != nil {
		return err
	}

	if err := j.client.Del(ctx, generateTokenValidityCacheKey(userId, tokenKey)).Err(); err != nil {
		// Dont return just log
		j.logger.Error(err)
	}

	return nil
}
