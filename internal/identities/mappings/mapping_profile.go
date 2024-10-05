package mappings

import (
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/models"
	userDtos "github.com/tguankheng016/golang-ecommerce-monolith/internal/identities/users/dtos"
	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/mapper"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.User, *userDtos.UserDto]()
	if err != nil {
		return err
	}
	return err
}
