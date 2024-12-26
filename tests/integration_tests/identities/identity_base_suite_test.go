package identities

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tguankheng016/commerce-mono/tests/integration_tests/shared"
)

type IdentityTestSuite struct {
	shared.AppTestSuite
}

func TestIdentitySuite(t *testing.T) {
	suite.Run(t, new(IdentityTestSuite))
}
