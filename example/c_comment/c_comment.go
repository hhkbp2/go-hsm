package c_comment

import hsm "github.com/hhkbp2/go-hsm"

// NewWorld() setup HSM and all states for this C Comment Parsing example.
func NewWorld() *CCommentHSM {
	// setup the state machine with proper hierarchy
	top := hsm.NewTop()
	initial := hsm.NewInitial(top, StateCodeID)
	NewCodeState(top)
	NewSlashState(top)
	NewCommentState(top)
	NewStarState(top)
	sm := NewCCommentHSM(top, initial)
	// don't forget to initialize the hsm and all states.
	sm.Init()
	return sm
}
