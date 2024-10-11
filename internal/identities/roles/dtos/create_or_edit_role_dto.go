package dtos

type CreateOrEditRoleDto struct {
	Id                 int64    `json:"id"`
	Name               string   `json:"name" validate:"required"`
	GrantedPermissions []string `json:"grantedPermissions"`
} // @name CreateOrEditRoleDto

type CreateRoleDto struct {
	*CreateOrEditRoleDto
} // @name CreateRoleDto

type EditRoleDto struct {
	*CreateOrEditRoleDto
} // @name EditRoleDto
