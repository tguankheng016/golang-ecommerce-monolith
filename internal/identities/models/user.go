package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"

	"gorm.io/gorm"
)

// User model
type User struct {
	Id            int64          `json:"id" gorm:"primaryKey"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	UserName      string         `json:"userName"`
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	SecurityStamp uuid.UUID      `json:"securityStamp"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"default:current_timestamp"`
	CreatedBy     sql.NullInt64  `json:"createdBy"`
	UpdatedAt     sql.NullTime   `json:"updatedAt"`
	UpdatedBy     sql.NullInt64  `json:"updatedBy"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
	Roles         []Role         `gorm:"many2many:user_roles;"`
}
