package dtos

type UserLoginInfoDto struct {
	Id        int64  `json:"id"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
