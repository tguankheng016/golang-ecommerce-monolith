package dtos

type PermissionDto struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsGranted   bool   `json:"isGranted"`
}

type PermissionGroupDto struct {
	GroupName   string          `json:"groupName"`
	Permissions []PermissionDto `json:"permissions"`
}
