package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Role model
type Role struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createdAt" gorm:"default:current_timestamp"`
	CreatedBy sql.NullInt64  `json:"createdBy"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
	UpdatedBy sql.NullInt64  `json:"updatedBy"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
