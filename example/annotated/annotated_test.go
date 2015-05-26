package annotated

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorld(t *testing.T) {
	sm := NewWorld()
	// Initial state is S11 state.
	assert.Equal(t, StateS11ID, sm.CurrentStateID())
	// Feed the state machine with input from event a to h.
	// check the state transfer done well when each event dispatched.
	// To digest how the state transfer process happens, we could also
	// trace every transfer action on the command line output.
	sm.Dispatch(NewEvent(EventA))
	assert.Equal(t, StateS11ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventB))
	assert.Equal(t, StateS11ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventC))
	assert.Equal(t, StateS211ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventD))
	assert.Equal(t, StateS211ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventE))
	assert.Equal(t, StateS211ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventF))
	assert.Equal(t, StateS11ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventG))
	assert.Equal(t, StateS211ID, sm.CurrentStateID())
	sm.Dispatch(NewEvent(EventH))
	assert.Equal(t, StateS211ID, sm.CurrentStateID())
}
