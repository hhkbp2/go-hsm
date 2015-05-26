package annotated

import hsm "github.com/hhkbp2/go-hsm"

func NewWorld() *AnnotatedHSM {
	top := hsm.NewTop()
	initial := hsm.NewInitial(top, StateS0ID)
	s0 := NewS0State(top)
	s1 := NewS1State(s0)
	NewS11State(s1)
	s2 := NewS2State(s0)
	s21 := NewS21State(s2)
	NewS211State(s21)
	sm := NewAnnotatedHSM(top, initial)
	sm.Init()
	return sm
}
