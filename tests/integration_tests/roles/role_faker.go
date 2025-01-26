package roles

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/tguankheng016/commerce-mono/internal/roles/models"
)

func GetFakeRole() *models.Role {
	role := models.Role{
		Id:   0,
		Name: gofakeit.BeerName(),
	}

	return &role
}
