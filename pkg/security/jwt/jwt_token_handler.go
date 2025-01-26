package jwt

import (
	"context"
	"errors"
	"strconv"

	jwtGo "github.com/golang-jwt/jwt/v5"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

type IJwtTokenHandler interface {
	ValidateToken(ctx context.Context, token string, tokenType TokenType) (int64, jwtGo.MapClaims, error)
}

type jwtTokenHandler struct {
	secretKey         string
	issuer            string
	audience          string
	stampValidator    IJwtSecurityStampValidator
	tokenKeyValidator IJwtTokenKeyValidator
}

func NewTokenHandler(authOptions *AuthOptions, stampValidator IJwtSecurityStampValidator, tokenKeyValidator IJwtTokenKeyValidator) IJwtTokenHandler {
	return &jwtTokenHandler{
		secretKey:         authOptions.SecretKey,
		issuer:            authOptions.Issuer,
		audience:          authOptions.Audience,
		stampValidator:    stampValidator,
		tokenKeyValidator: tokenKeyValidator,
	}
}

func (j *jwtTokenHandler) ValidateToken(ctx context.Context, tokenString string, tokenType TokenType) (int64, jwtGo.MapClaims, error) {
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
			return 0, nil, errors.New("invalid token type")
		}

		// token is valid and has not expired
		iss := token.Header["iss"]
		if iss != j.issuer {
			// handle invalid issuer
			return 0, nil, errors.New("invalid token issuer")
		}

		aud := token.Header["aud"]
		if aud != j.audience {
			// handle invalid audience
			return 0, nil, errors.New("invalid token audience")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			// handle error
			return 0, nil, errors.New("invalid sub")
		}

		userId, err := strconv.ParseInt(sub, 10, 64)

		if err != nil {
			return 0, nil, err
		}

		if err := j.stampValidator.ValidateTokenWithStamp(ctx, userId, claims); err != nil {
			return 0, nil, err
		}

		if err := j.tokenKeyValidator.ValidateTokenWithTokenKey(ctx, userId, claims); err != nil {
			return 0, nil, err
		}

		return userId, claims, nil
	}

	return 0, nil, errors.New("invalid token")
}
