package tests

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type YadosTestSuite struct {
	suite.Suite
}

func (suite *YadosTestSuite) SetupTest() {
}

func (suite *YadosTestSuite) Cleanup() {
}

func (suite *YadosTestSuite) TearDownTest() {
}

func TestYadosTestSuite(t *testing.T) {
	suite.Run(t, new(YadosTestSuite))
}
