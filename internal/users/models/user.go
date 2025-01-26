package models

import (
	"github.com/gofrs/uuid"
	"github.com/tguankheng016/commerce-mono/pkg/core/domain"
)

type User struct {
	Id                 int64
	FirstName          string
	LastName           string
	UserName           string
	NormalizedUserName string
	Email              string
	NormalizedEmail    string
	PasswordHash       string
	SecurityStamp      uuid.UUID
	domain.FullAuditedEntity
}

type UserRole struct {
	UserId int64
	RoleId int64
	domain.CreationAuditedEntity
}
