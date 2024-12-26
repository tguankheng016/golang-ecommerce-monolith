package models

import (
	"database/sql"
	"time"
)

type UserRolePermission struct {
	Id        uint
	Name      string
	UserId    sql.NullInt64
	RoleId    sql.NullInt64
	IsGranted bool
	CreatedAt time.Time
}
