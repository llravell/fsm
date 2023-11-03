package fsm

import (
	"fmt"
)

type Event[S comparable] struct {
	Name string
	From []S
	To   S
}

type Events[S comparable] []Event[S]

type UnknownEventError struct {
	event string
}

func (e UnknownEventError) Error() string {
	return fmt.Sprintf("[fsm] event \"%s\" has not been defined", e.event)
}

type InvalidStateError[S comparable] struct {
	from S
	to   S
}

func (e InvalidStateError[S]) Error() string {
	return fmt.Sprintf("[fsm] can not move from %v to %v", e.from, e.to)
}

type eventsMap[S comparable] map[string]S

type transactionsMap[S comparable] map[S]*set[S]

type FSM[S comparable] struct {
	state S
	em    eventsMap[S]
	tm    transactionsMap[S]
}

func NewFSM[S comparable](initial S, events Events[S]) FSM[S] {
	return FSM[S]{
		state: initial,
		em:    buildEventsMap(events),
		tm:    buildTransactionsMap(events),
	}
}

func (fsm FSM[S]) GetState() S {
	return fsm.state
}

func (fsm FSM[S]) GetAvailableStates() []S {
	availableStates, has := fsm.tm[fsm.GetState()]

	if !has {
		return make([]S, 0)
	}

	return availableStates.Keys()
}

func (fsm FSM[S]) CanMoveTo(s S) bool {
	availableStates, has := fsm.tm[fsm.GetState()]
	if !has {
		return false
	}

	return availableStates.Has(s)
}

func (fsm *FSM[S]) Event(event string) error {
	targetState, has := fsm.em[event]

	if !has {
		return &UnknownEventError{event}
	}

	if !fsm.CanMoveTo(targetState) {
		return &InvalidStateError[S]{fsm.GetState(), targetState}
	}

	fsm.state = targetState
	return nil
}

func buildTransactionsMap[S comparable](events Events[S]) transactionsMap[S] {
	tm := make(transactionsMap[S])

	for _, event := range events {
		for _, from := range event.From {
			_, has := tm[from]

			if !has {
				set := newSet[S]()
				tm[from] = &set
			}

			tm[from].Add(event.To)
		}
	}

	return tm
}

func buildEventsMap[S comparable](events Events[S]) eventsMap[S] {
	em := make(eventsMap[S])

	for _, event := range events {
		em[event.Name] = event.To
	}

	return em
}
