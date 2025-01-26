package jwt

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tguankheng016/commerce-mono/pkg/http"
	"github.com/tguankheng016/commerce-mono/pkg/logging"
	"go.uber.org/zap"
)

func SetupJwtAuthentication(api huma.API, tokenHandler IJwtTokenHandler) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Parse and verify jwt access token
		authToken, ok := bearerAuth(ctx)
		if !ok {
			logging.Logger.Warn("no access token found in the authorization header")
			next(ctx)
			return
		}

		context := ctx.Context()

		// Validate jwt access token
		userId, claims, err := tokenHandler.ValidateToken(context, authToken, AccessToken)
		if err != nil {
			logging.Logger.Error("validate jwt access token error: ", zap.Error(err))
			next(ctx)
			return
		}

		ctx = http.SetCurrentUser(ctx, userId)
		ctx = http.SetCurrentUserClaims(ctx, claims)

		next(ctx)
	}
}

// BearerAuth parse bearer token
func bearerAuth(ctx huma.Context) (string, bool) {
	auth := ctx.Header("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	return token, token != ""
}
