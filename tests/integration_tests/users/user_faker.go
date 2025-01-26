package users

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofrs/uuid"
	"github.com/tguankheng016/commerce-mono/internal/users/models"
)

func GetFakeUser() *models.User {
	securityStamp, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	user := models.User{
		Id:            0,
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		UserName:      gofakeit.Username(),
		Email:         gofakeit.Email(),
		SecurityStamp: securityStamp,
	}

	return &user
}
