package models

import (
	"database/sql"
	"time"
)

type UserRolePermission struct {
	Id        uint      `gorm:"primary_key"`
	Name      string    `gorm:"type:varchar(256);not null"`
	UserId    int64     `gorm:"index"`
	RoleId    int64     `gorm:"index"`
	IsGranted bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt sql.NullTime
}
