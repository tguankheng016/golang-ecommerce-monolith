package dtos

type CreateOrEditUserDto struct {
	Id        int64   `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	UserName  string  `json:"userName"`
	Email     string  `json:"email"`
	Password  string  `json:"password" copier:"-"`
	RoleIds   []int64 `json:"roleIds"`
}

type CreateUserDto struct {
	*CreateOrEditUserDto
} // @name CreateUserDto

type EditUserDto struct {
	*CreateOrEditUserDto
}
