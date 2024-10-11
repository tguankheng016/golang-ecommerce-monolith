package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoServer "github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/http/echo"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/jwt"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/logger"
)

// SetupAuthenticate is a middleware that try to authenticate the user with the given
// `tokenValidator` using the Authorization header.
//
// If the header is not present or the token is invalid, the middleware will call the
// next handler as usual.
//
// If the token is valid, the middleware will set the current user id into the request
// context, and then call the next handler.
//
// The `config.Skipper` can be used to skip the authentication process.
func SetupAuthenticate(skipper echoMiddleware.Skipper, tokenValidator jwt.IJwtTokenValidator, logger logger.ILogger) echo.MiddlewareFunc {
	// Defaults
	if skipper == nil {
		skipper = echoMiddleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			// Parse and verify jwt access token
			authToken, ok := bearerAuth(c.Request())
			if !ok {
				logger.Warn("No access token found in the Authorization header")
				return next(c)
			}

			ctx := c.Request().Context()

			// Validate jwt access token
			userId, claims, err := tokenValidator.ValidateToken(ctx, authToken, jwt.AccessToken)
			if err != nil {
				logger.Error(err)
				return next(c)
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
