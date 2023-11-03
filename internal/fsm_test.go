package fsm

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type State string

const (
	Inactive State = "inactive"
	Active   State = "active"
	Frozen   State = "frozen"
)

type FSMTestSuite struct {
	suite.Suite
	fsm *FSM[State]
}

func (suite *FSMTestSuite) SetupTest() {
	fsm := NewFSM[State](Inactive, Events[State]{
		{Name: "activate", From: []State{Inactive}, To: Active},
		{Name: "dectivate", From: []State{Active}, To: Inactive},
		{Name: "freeze", From: []State{Active}, To: Frozen},
	})

	suite.fsm = &fsm
}

func (suite *FSMTestSuite) TestEventPositive() {
	err := suite.fsm.Event("activate")
	suite.Nil(err)
	suite.Equal(Active, suite.fsm.GetState())

	err = suite.fsm.Event("dectivate")
	suite.Nil(err)
	suite.Equal(Inactive, suite.fsm.GetState())
}

func (suite *FSMTestSuite) TestEventInvalid() {
	var invalidStateError InvalidStateError[State]
	err := suite.fsm.Event("freeze")

	suite.ErrorAs(invalidStateError, err)
	suite.Equal(Inactive, suite.fsm.GetState())

	suite.fsm.Event("activate")
	err = suite.fsm.Event("freeze")

	suite.Nil(err)
	suite.Equal(Frozen, suite.fsm.GetState())

	err = suite.fsm.Event("activate")

	suite.ErrorAs(invalidStateError, err)
	suite.Equal(Frozen, suite.fsm.GetState())
}

func (suite *FSMTestSuite) TestEventUnknown() {
	var unknownEventError UnknownEventError
	err := suite.fsm.Event("unknown event")

	suite.ErrorAs(unknownEventError, err)
	suite.Equal(Inactive, suite.fsm.GetState())
}

func (suite *FSMTestSuite) TestCanMoveTo() {
	suite.True(suite.fsm.CanMoveTo(Active))
	suite.False(suite.fsm.CanMoveTo(Inactive))
	suite.False(suite.fsm.CanMoveTo(Frozen))

	suite.fsm.Event("activate")

	suite.False(suite.fsm.CanMoveTo(Active))
	suite.True(suite.fsm.CanMoveTo(Inactive))
	suite.True(suite.fsm.CanMoveTo(Frozen))

	suite.fsm.Event("freeze")

	suite.False(suite.fsm.CanMoveTo(Active))
	suite.False(suite.fsm.CanMoveTo(Inactive))
	suite.False(suite.fsm.CanMoveTo(Frozen))
}

func (suite *FSMTestSuite) TestGetAvailableStates() {
	suite.ElementsMatch(
		[]State{Active},
		suite.fsm.GetAvailableStates(),
	)

	suite.fsm.Event("activate")

	suite.ElementsMatch(
		[]State{Inactive, Frozen},
		suite.fsm.GetAvailableStates(),
	)

	suite.fsm.Event("freeze")

	suite.Empty(suite.fsm.GetAvailableStates())
}

func TestFSMTestSuite(t *testing.T) {
	suite.Run(t, new(FSMTestSuite))
}
