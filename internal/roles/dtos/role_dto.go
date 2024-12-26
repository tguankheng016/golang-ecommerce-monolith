package dtos

import "github.com/tguankheng016/commerce-mono/pkg/core/domain"

type RoleDto struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsStatic  string `json:"isStatic"`
	domain.AuditedEntityDto
}
