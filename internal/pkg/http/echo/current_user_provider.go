package echoserver

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	userIdkey         = "currentUser:userId"
	claimsKey         = "currentUser:claims"
	userPermissionKey = "currentUser:permissions"
)

func SetCurrentUser(c echo.Context, userId int64) {
	c.Set(userIdkey, userId)
}

func GetCurrentUser(c echo.Context) (int64, bool) {
	userId, ok := c.Get(userIdkey).(int64)
	return userId, ok
}

func SetCurrentUserClaims(c echo.Context, claims jwt.MapClaims) {
	c.Set(claimsKey, claims)
}

func GetCurrentUserClaims(c echo.Context) (jwt.MapClaims, bool) {
	claims, ok := c.Get(claimsKey).(jwt.MapClaims)
	return claims, ok
}

func SetCurrentUserPermissions(c echo.Context, permissions map[string]struct{}) {
	c.Set(userPermissionKey, permissions)
}

func GetCurrentUserPermissions(c echo.Context) (map[string]struct{}, bool) {
	permissions, ok := c.Get(userPermissionKey).(map[string]struct{})
	return permissions, ok
}
