package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	userIdkey         = "currentUser:userId"
	claimsKey         = "currentUser:claims"
	userPermissionKey = "currentUser:permissions"
)

func SetCurrentUser(ctx huma.Context, userId int64) huma.Context {
	ctx = huma.WithValue(ctx, userIdkey, userId)
	return ctx
}

func GetCurrentUser(ctx context.Context) (int64, bool) {
	userId, ok := ctx.Value(userIdkey).(int64)
	return userId, ok
}

func SetCurrentUserClaims(ctx huma.Context, claims jwt.MapClaims) huma.Context {
	ctx = huma.WithValue(ctx, claimsKey, claims)
	return ctx
}

func GetCurrentUserClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.MapClaims)
	return claims, ok
}

func SetCurrentUserPermissions(ctx huma.Context, permissions map[string]struct{}) huma.Context {
	ctx = huma.WithValue(ctx, userPermissionKey, permissions)
	return ctx
}

func GetCurrentUserPermissions(ctx context.Context) (map[string]struct{}, bool) {
	permissions, ok := ctx.Value(userPermissionKey).(map[string]struct{})
	return permissions, ok
}
