package domain

import "database/sql"

type CreationAuditedEntity struct {
	CreatedBy sql.NullInt64 `json:"createdBy"`
}

type UpdateAuditedEntity struct {
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
}

type DeleteAuditedEntity struct {
	DeletedBy sql.NullInt64 `json:"deletedBy"`
}

type FullAuditedEntity struct {
	*CreationAuditedEntity
	*UpdateAuditedEntity
	*DeleteAuditedEntity
}
