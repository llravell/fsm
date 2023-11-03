package fsm

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SetTestSuite struct {
	suite.Suite
	set set[string]
}

func (suite *SetTestSuite) SetupTest() {
	suite.set = newSet[string]()
}

func (suite *SetTestSuite) TestAdd() {
	suite.set.Add("a")
	suite.set.Add("b")

	suite.True(suite.set.Has("a"))
	suite.True(suite.set.Has("b"))
}

func (suite *SetTestSuite) TestDelete() {
	suite.set.Add("a")
	suite.True(suite.set.Has("a"))

	suite.set.Delete("a")
	suite.False(suite.set.Has("a"))
}

func (suite *SetTestSuite) TestSize() {
	suite.set.Add("a")
	suite.set.Add("b")
	suite.set.Add("c")
	suite.set.Add("c")
	suite.set.Add("c")

	suite.Equal(3, suite.set.Size())

	suite.set.Delete("c")
	suite.Equal(2, suite.set.Size())
}

func (suite *SetTestSuite) TestKeys() {
	suite.set.Add("a")
	suite.set.Add("b")
	suite.set.Add("c")
	suite.set.Add("c")
	suite.set.Add("c")

	suite.ElementsMatch([]string{"a", "b", "c"}, suite.set.Keys())
}

func TestSetTestSuite(t *testing.T) {
	suite.Run(t, new(SetTestSuite))
}
