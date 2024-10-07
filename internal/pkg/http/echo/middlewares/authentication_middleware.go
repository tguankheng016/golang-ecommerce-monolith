package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
)

func ValidateToken(validator jwt.IJwtTokenValidator) echo.MiddlewareFunc {
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
				return echo.NewHTTPError(http.StatusUnauthorized, errors.New("parse token error"))
			}

			// Validate jwt access token
			userId, _, err := validator.ValidateToken(auth, jwt.AccessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			c.Set("userId", userId)

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
