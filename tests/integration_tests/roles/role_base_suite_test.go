package roles

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tguankheng016/commerce-mono/tests/integration_tests/shared"
)

type RoleTestSuite struct {
	shared.AppTestSuite
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}

func (s *RoleTestSuite) ResetRoles() {
	query := "DELETE FROM roles WHERE id > 2"

	if _, err := s.Pool.Exec(s.Ctx, query); err != nil {
		panic(err)
	}
}

func (s *RoleTestSuite) ResetRolePermissions() {
	query := "DELETE FROM user_role_permissions WHERE role_id is not null"

	if _, err := s.Pool.Exec(s.Ctx, query); err != nil {
		panic(err)
	}
}
