package jwt

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
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
	GenerateAccessToken(user *models.User, refreshTokenKey string) (string, error)
}

type jwtTokenGenerator struct {
	secretKey string
	issuer    string
	audience  string
	db        *gorm.DB
}

func NewJwtTokenGenerator(db *gorm.DB, authOptions *AuthOptions) IJwtTokenGenerator {
	return &jwtTokenGenerator{
		secretKey: authOptions.SecretKey,
		issuer:    authOptions.Issuer,
		audience:  authOptions.Audience,
		db:        db,
	}
}

func (j *jwtTokenGenerator) GenerateAccessToken(user *models.User, refreshTokenKey string) (string, error) {
	claims, err := j.createJwtClaims(user, AccessToken, refreshTokenKey)

	if err != nil {
		return "", err
	}

	return j.createToken(claims)
}

func (j *jwtTokenGenerator) createToken(claims jwtGo.MapClaims) (string, error) {
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token.Header["iss"] = j.issuer
	token.Header["aud"] = j.audience

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtTokenGenerator) createJwtClaims(user *models.User, tokenType TokenType, refreshTokenKey string) (jwtGo.MapClaims, error) {
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

	return claims, nil
}
