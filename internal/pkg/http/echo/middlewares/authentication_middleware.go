package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
)

// ValidateToken validates the JWT token in the Authorization header. If the token is invalid, it returns a 401 Unauthorized response.
func ValidateToken(validator jwt.IJwtTokenValidator) echo.MiddlewareFunc {
	return validateToken(validator, false)
}

// TryValidateToken validates the JWT token in the Authorization header but does not throw an error if the token is invalid.
func TryValidateToken(validator jwt.IJwtTokenValidator) echo.MiddlewareFunc {
	return validateToken(validator, true)
}

func validateToken(validator jwt.IJwtTokenValidator, canSkipError bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Ignore check authentication in test
			env := os.Getenv("APP_ENV")
			if env == "test" {
				return next(c)
			}

			// Parse and verify jwt access token
			auth, ok := bearerAuth(c.Request())
			if !ok {
				if !canSkipError {
					return echo.NewHTTPError(http.StatusUnauthorized, errors.New("Invalid access token"))
				} else {
					return next(c)
				}
			}

			ctx := c.Request().Context()

			// Validate jwt access token
			userId, claims, err := validator.ValidateToken(ctx, auth, jwt.AccessToken)
			if err != nil {
				if !canSkipError {
					return echo.NewHTTPError(http.StatusUnauthorized, err)
				} else {
					return next(c)
				}
			}

			echoServer.SetCurrentUser(c, userId)
			echoServer.SetCurrentUserClaims(c, claims)

			return next(c)
		}
	}
}

// BearerAuth parse bearer token
func bearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access_token")
	}
	return token, token != ""
}
