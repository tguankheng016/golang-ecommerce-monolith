package users

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tguankheng016/commerce-mono/tests/integration_tests/shared"
)

type UserTestSuite struct {
	shared.AppTestSuite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) ResetUsers() {
	s.ResetUserRoles()
	s.ResetUserPermissions()

	query := "DELETE FROM users WHERE id > 2"

	if _, err := s.Pool.Exec(s.Ctx, query); err != nil {
		panic(err)
	}
}

func (s *UserTestSuite) ResetUserRoles() {
	query := "DELETE FROM user_roles WHERE user_id > 2"

	if _, err := s.Pool.Exec(s.Ctx, query); err != nil {
		panic(err)
	}
}

func (s *UserTestSuite) ResetUserPermissions() {
	query := "DELETE FROM user_role_permissions WHERE user_id is not null AND user_id > 2"

	if _, err := s.Pool.Exec(s.Ctx, query); err != nil {
		panic(err)
	}
}
