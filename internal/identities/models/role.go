package models

import (
	"database/sql"
	"time"

	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/core/domain"
	"gorm.io/gorm"
)

// Role model
type Role struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(256);not null"`
	CreatedAt time.Time      `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	*domain.FullAuditedEntity
}
