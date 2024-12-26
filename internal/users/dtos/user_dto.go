package dtos

import "github.com/tguankheng016/commerce-mono/pkg/core/domain"

type UserDto struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	domain.AuditedEntityDto
}
