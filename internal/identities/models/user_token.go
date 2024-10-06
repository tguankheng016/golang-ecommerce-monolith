package models

import "time"

// User Token model
type UserToken struct {
	Id             int64     `gorm:"primarykey"`
	UserId         int64     `gorm:"column:user_id;index"`
	TokenKey       string    `gorm:"column:token_key"`
	ExpirationTime time.Time `gorm:"column:expiration_time"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
