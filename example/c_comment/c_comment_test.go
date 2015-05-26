package c_comment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestWorld() is a use case of this example state machine.
func TestWorld(t *testing.T) {
	sm := NewWorld()
	// Initial state is code state.
	assert.Equal(t, StateCodeID, sm.CurrentStateID())
	// Feed the state machine with input:
	// a/=/*c/*d**/b
	// When all event dispatched, the state machine should record code like:
	// a/=b
	sm.Dispatch(NewCCommentCharEvent('a'))
	// still in code state
	assert.Equal(t, StateCodeID, sm.CurrentStateID())
	// transfer to slash state
	sm.Dispatch(NewCCommentSlashEvent())
	assert.Equal(t, StateSlashID, sm.CurrentStateID())
	// transfer to code state
	sm.Dispatch(NewCCommentCharEvent('='))
	assert.Equal(t, StateCodeID, sm.CurrentStateID())
	// transfer to slash state, again
	sm.Dispatch(NewCCommentSlashEvent())
	assert.Equal(t, StateSlashID, sm.CurrentStateID())
	// transfer to comment state
	sm.Dispatch(NewCCommentStarEvent())
	assert.Equal(t, StateCommentID, sm.CurrentStateID())
	sm.Dispatch(NewCCommentCharEvent('c'))
	sm.Dispatch(NewCCommentSlashEvent())
	assert.Equal(t, StateCommentID, sm.CurrentStateID())
	// transfer to star state
	sm.Dispatch(NewCCommentStarEvent())
	assert.Equal(t, StateStarID, sm.CurrentStateID())
	// transfer to comment state, again
	sm.Dispatch(NewCCommentCharEvent('d'))
	assert.Equal(t, StateCommentID, sm.CurrentStateID())
	// transfer to star state
	sm.Dispatch(NewCCommentStarEvent())
	assert.Equal(t, StateStarID, sm.CurrentStateID())
	sm.Dispatch(NewCCommentStarEvent())
	assert.Equal(t, StateStarID, sm.CurrentStateID())
	// transfer to code state
	sm.Dispatch(NewCCommentSlashEvent())
	assert.Equal(t, StateCodeID, sm.CurrentStateID())
	sm.Dispatch(NewCCommentCharEvent('b'))

	// traverse the code chars and check the record code chars
	codeChars := make([]byte, 0)
	f := func(value interface{}) interface{} {
		b, ok := value.(byte)
		assert.Equal(t, true, ok)
		codeChars = append(codeChars, b)
		return nil
	}
	sm.TraverseCode(f)
	assert.Equal(t, "a/=b", string(codeChars))
}
