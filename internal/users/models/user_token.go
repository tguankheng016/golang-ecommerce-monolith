package models

import "time"

type UserToken struct {
	Id             int64
	UserId         int64
	TokenKey       string
	ExpirationTime time.Time
}
