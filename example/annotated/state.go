package annotated

import hsm "github.com/hhkbp2/go-hsm"
import "log"

const (
	StateS0ID   string = "s0"
	StateS1ID          = "s1"
	StateS11ID         = "s11"
	StateS2ID          = "s2"
	StateS21ID         = "s21"
	StateS211ID        = "s211"
)

func Logln(v ...interface{}) {
	log.Println(v...)
}

type VerboseStateHead struct {
	*hsm.StateHead
	ID string
}

func NewVerboseStateHead(super hsm.State) *VerboseStateHead {
	return &VerboseStateHead{
		StateHead: hsm.NewStateHead(super),
	}
}

func (self *VerboseStateHead) Init(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Init")
	return nil
}

func (self *VerboseStateHead) Entry(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Entry")
	return nil
}

func (self *VerboseStateHead) Exit(hsm hsm.HSM, event hsm.Event) (state hsm.State) {
	Logln(self.ID, "- Exit")
	return nil
}

/*
 * Feel sorry to write these repeatly due to the lack of
 * meta programming ability in Golang.
 */

type S0State struct {
	*VerboseStateHead
}

func NewS0State(super hsm.State) *S0State {
	object := &S0State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (_ *S0State) ID() string {
	return StateS0ID
}

func (self *S0State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateS1ID)
	return nil
}

func (self *S0State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventE:
		sm.QTran(StateS211ID)
		return nil
	}
	return self.Super()
}

type S1State struct {
	*VerboseStateHead
}

func NewS1State(super hsm.State) *S1State {
	object := &S1State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (_ *S1State) ID() string {
	return StateS1ID
}

func (self *S1State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateS11ID)
	return nil
}

func (self *S1State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventA:
		sm.QTran(StateS1ID)
		return nil
	case EventB:
		sm.QTran(StateS11ID)
		return nil
	case EventC:
		sm.QTran(StateS2ID)
		return nil
	case EventD:
		sm.QTran(StateS0ID)
		return nil
	case EventF:
		sm.QTran(StateS211ID)
		return nil
	}
	return self.Super()
}

type S11State struct {
	*VerboseStateHead
}

func NewS11State(super hsm.State) *S11State {
	object := &S11State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (_ *S11State) ID() string {
	return StateS11ID
}

func (self *S11State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}

func (self *S11State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventG:
		sm.QTran(StateS211ID)
		return nil
	case EventH:
		annotatedHSM, ok := sm.(*AnnotatedHSM)
		hsm.AssertTrue(ok)
		if annotatedHSM.GetFoo() {
			annotatedHSM.SetFoo(false)
			return nil
		}
	}
	return self.Super()
}

type S2State struct {
	*VerboseStateHead
}

func NewS2State(super hsm.State) *S2State {
	object := &S2State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (_ *S2State) ID() string {
	return StateS2ID
}

func (self *S2State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateS21ID)
	return nil
}

func (self *S2State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventC:
		sm.QTran(StateS1ID)
		return nil
	case EventF:
		sm.QTran(StateS11ID)
		return nil
	}
	return self.Super()
}

type S21State struct {
	*VerboseStateHead
}

func NewS21State(super hsm.State) *S21State {
	object := &S21State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (_ *S21State) ID() string {
	return StateS21ID
}

func (self *S21State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	sm.QInit(StateS211ID)
	return nil
}

func (self *S21State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventB:
		sm.QTran(StateS211ID)
		return nil
	case EventH:
		annotatedHSM, ok := sm.(*AnnotatedHSM)
		hsm.AssertTrue(ok)
		if !annotatedHSM.GetFoo() {
			annotatedHSM.SetFoo(true)
			sm.QTran(StateS21ID)
			return nil
		}
	}
	return self.Super()
}

type S211State struct {
	*VerboseStateHead
}

func NewS211State(super hsm.State) *S211State {
	object := &S211State{
		NewVerboseStateHead(super),
	}
	object.VerboseStateHead.ID = object.ID()
	super.AddChild(object)
	return object
}

func (self *S211State) ID() string {
	return StateS211ID
}

func (self *S211State) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	self.VerboseStateHead.Init(sm, event)
	return self.Super()
}

func (self *S211State) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle e =", PrintEvent(event.Type()))
	switch event.Type() {
	case EventD:
		sm.QTran(StateS21ID)
		return nil
	case EventG:
		sm.QTran(StateS0ID)
		return nil
	}
	return self.Super()
}
