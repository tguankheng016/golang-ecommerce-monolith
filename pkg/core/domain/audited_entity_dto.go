package domain

import (
	"database/sql"
	"time"
)

type CreationAuditedEntityDto struct {
	CreatedAt time.Time     `json:"createdAt"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
}

type AuditedEntityDto struct {
	UpdatedAt sql.NullTime  `json:"updatedAt"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
	CreationAuditedEntityDto
}
