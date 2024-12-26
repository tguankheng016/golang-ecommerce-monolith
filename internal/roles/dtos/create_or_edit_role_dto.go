package dtos

type CreateOrEditRoleDto struct {
	Id                 *int64   `json:"id"`
	Name               string   `json:"name"`
	IsDefault          bool     `json:"isDefault"`
	GrantedPermissions []string `json:"grantedPermissions"`
}

type CreateRoleDto struct {
	CreateOrEditRoleDto
}

type EditRoleDto struct {
	CreateOrEditRoleDto
}
