package dtos

type CreateOrEditUserDto struct {
	Id        int64   `json:"id"`
	FirstName string  `json:"firstName" validate:"min=3,max=64,required"`
	LastName  string  `json:"lastName" validate:"min=3,max=64,required"`
	UserName  string  `json:"userName" validate:"min=8,max=256,required"`
	Email     string  `json:"email" validate:"email,max=256,required"`
	Password  string  `json:"password" copier:"-"`
	RoleIds   []int64 `json:"roleIds"`
} // @name CreateOrEditUserDto

type CreateUserDto struct {
	*CreateOrEditUserDto
} // @name CreateUserDto

type EditUserDto struct {
	*CreateOrEditUserDto
} // @name EditUserDto
